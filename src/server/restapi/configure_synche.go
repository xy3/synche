// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/auth"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/handlers"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/tokens"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/users"
	"net/http"
)

// //go:generate swagger generate server --target ../../server --name Synche --spec ../spec/synche-server-api.yaml --principal gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema.User --flag-strategy=pflag --exclude-main

func configureFlags(api *operations.SyncheAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SyncheAPI) http.Handler {
	api.ServeError = errors.ServeError
	api.Logger = log.Infof
	api.UseSwaggerUI()
	api.JSONConsumer = runtime.JSONConsumer()
	api.MultipartformConsumer = runtime.DiscardConsumer
	api.JSONProducer = runtime.JSONProducer()
	api.BinProducer = runtime.ByteStreamProducer()

	authService := auth.Service{
		SecretKey:       c.Config.Server.SecretKey,
		Issuer:          "synche.auth.service",
		ExpirationHours: 24,
	}

	// Applies when the "Bearer" header is set
	api.AccessTokenAuth = func(token string) (*schema.User, error) {
		if cachedUser, found := repo.TokenToUserCache.Get(token); found {
			return cachedUser.(*schema.User), nil
		}

		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			return nil, fmt.Errorf("could not validate access token")
		}

		user, err := repo.GetUserByEmail(claims.Email, database.DB)
		if err != nil {
			return nil, fmt.Errorf("user credentials not found")
		}

		repo.TokenToUserCache.SetDefault(token, user)

		return user, nil
	}

	api.RefreshTokenAuth = func(token string) (*schema.User, error) {
		claims, err := authService.ValidateRefreshToken(token)
		if err != nil {
			return nil, errors.New(400, "token could not be validated")
		}
		user, err := repo.GetUserByEmail(claims.Email, database.DB)
		if err != nil {
			return nil, errors.New(404, "user credentials not found")
		}

		actualCustomKey := authService.GenerateCustomKey(user.Email, user.TokenHash)
		if claims.CustomKey != actualCustomKey {
			return nil, errors.New(400, "invalid token, authentication failed")
		}

		return user, nil
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// api.APIAuthorizer = security.Authorized()
	// You may change here the memory limit for this multipart form parser. Currently 50 MB.
	transfer.UploadChunkMaxParseMemory = 50 << 20

	// User handlers
	api.UsersLoginHandler = users.LoginHandlerFunc(func(params users.LoginParams) middleware.Responder {
		return handlers.Login(params, authService)
	})
	api.UsersRegisterHandler = users.RegisterHandlerFunc(handlers.Register)
	api.UsersProfileHandler = users.ProfileHandlerFunc(handlers.Profile)
	api.UsersDeleteUserHandler = users.DeleteUserHandlerFunc(handlers.DeleteUser)

	// File handlers
	api.FilesDeleteFileHandler = files.DeleteFileHandlerFunc(handlers.DeleteFileID)
	api.FilesGetFileInfoHandler = files.GetFileInfoHandlerFunc(handlers.FileInfo)
	api.FilesUpdateFileByIDHandler = files.UpdateFileByIDHandlerFunc(handlers.UpdateFileByID)
	api.FilesDeleteFilepathHandler = files.DeleteFilepathHandlerFunc(handlers.DeleteFilePath)
	api.FilesUpdateFileByPathHandler = files.UpdateFileByPathHandlerFunc(handlers.UpdateFileByPath)
	api.FilesGetFilePathInfoHandler = files.GetFilePathInfoHandlerFunc(handlers.FilePathInfo)

	// Directory handlers
	api.FilesListDirectoryHandler = files.ListDirectoryHandlerFunc(handlers.ListDirectory)
	api.FilesListHomeDirectoryHandler = files.ListHomeDirectoryHandlerFunc(handlers.ListHomeDirectory)
	api.FilesCreateDirectoryHandler = files.CreateDirectoryHandlerFunc(handlers.CreateDirectory)
	api.FilesDeleteDirectoryHandler = files.DeleteDirectoryHandlerFunc(handlers.DeleteDirectory)

	// Transfer handlers
	api.TransferDownloadFileHandler = transfer.DownloadFileHandlerFunc(handlers.DownloadFile)
	api.TransferNewUploadHandler = transfer.NewUploadHandlerFunc(handlers.NewUpload)
	api.TransferUploadChunkHandler = transfer.UploadChunkHandlerFunc(handlers.UploadChunk)
	api.TransferFinishUploadHandler = transfer.FinishUploadHandlerFunc(handlers.FinishUpload)

	api.TokensRefreshTokenHandler = tokens.RefreshTokenHandlerFunc(func(
		params tokens.RefreshTokenParams,
		user *schema.User,
	) middleware.Responder {
		token, err := authService.GenerateAccessToken(user.ID, user.Email, user.Name, user.Picture, user.Role)
		if err != nil {
			return tokens.NewRefreshTokenDefault(500).WithPayload("could not generate a new access token")
		}
		return tokens.NewRefreshTokenOK().WithPayload(&models.AccessToken{AccessToken: token})
	})
	// 	============= End Route Handlers =============

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {}
	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handlerCors := cors.AllowAll().Handler
	return handlerCors(handler)
}
