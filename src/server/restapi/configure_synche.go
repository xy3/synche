// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/auth"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/handlers"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/tokens"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/users"
	"net/http"

	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
)

//go:generate swagger generate server --target ../../server --name Synche --spec ../spec/synche-server-api.yaml --principal models.User --exclude-main

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

	err := data.InitSyncheData()
	if err != nil {
		log.WithError(err).Fatal("Failed to start Synche Data requirements")
	}

	authService := auth.Service{
		SecretKey:       "CHANGE_THIS_SECRET_KEY", // TODO: This should be configurable
		Issuer:          "synche.auth.service",
		ExpirationHours: 24,
	}

	// Applies when the "Bearer" header is set
	api.AccessTokenAuth = func(token string) (*schema.User, error) {
		log.Info(token)
		claims, err := authService.ValidateAccessToken(token)
		if err != nil {
			return nil, errors.New(400, err.Error())
		}

		user, err := repo.GetUserByEmail(claims.Email)
		if err != nil {
			return nil, errors.New(404, "user credentials not found")
		}
		return user, nil
	}

	api.RefreshTokenAuth = func(token string) (*schema.User, error) {
		claims, err := authService.ValidateRefreshToken(token)
		if err != nil {
			return nil, err
		}
		user, err := repo.GetUserByEmail(claims.Email)
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

	api.UsersLoginHandler = users.LoginHandlerFunc(func(params users.LoginParams) middleware.Responder {
		return handlers.Login(params, authService)
	})
	api.UsersRegisterHandler = users.RegisterHandlerFunc(handlers.Register)
	api.UsersProfileHandler = users.ProfileHandlerFunc(handlers.Profile)

	api.FilesDeleteFileHandler = files.DeleteFileHandlerFunc(handlers.DeleteFile)
	api.FilesGetFileInfoHandler = files.GetFileInfoHandlerFunc(handlers.FileInfo)
	api.FilesListHandler = files.ListHandlerFunc(handlers.ListFiles)

	api.TransferDownloadFileHandler = transfer.DownloadFileHandlerFunc(handlers.DownloadFile)
	api.TransferNewUploadHandler = transfer.NewUploadHandlerFunc(handlers.NewUpload)
	api.TransferUploadChunkHandler = transfer.UploadChunkHandlerFunc(handlers.UploadChunk)

	api.TokensRefreshTokenHandler = tokens.RefreshTokenHandlerFunc(func(
		params tokens.RefreshTokenParams,
		user *schema.User,
	) middleware.Responder {
		token, err := authService.GenerateAccessToken(user.Email)
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
	return handler
}
