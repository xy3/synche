// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/handlers"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
	"net/http"
)

//go:generate swagger generate server --target ../../server --name Synche --spec ../api/openapi-spec/synche-server-api.yaml --principal models.Message --flag-strategy=pflag --exclude-main

func configureFlags(api *operations.SyncheAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func OverrideFlags() {
	port = c.Config.Server.Port
}

func configureAPI(api *operations.SyncheAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Create data with chunk table and connection_request table if they don't exist
	err := data.CreateDatabase(c.Config.Database)
	if err != nil {
		log.WithError(err).Fatal("Database creation failed")
	}

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	api.Logger = log.Infof

	api.UseSwaggerUI()
	// To use redoc as the UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// transfer.UploadFileMaxParseMemory = 32 << 20

	// ============= Start Route Handlers =============
	// TODO: Implement listing functionality
	if api.TransferListFilesHandler == nil {
		api.TransferListFilesHandler = transfer.ListFilesHandlerFunc(func(params transfer.ListFilesParams) middleware.Responder {
			return middleware.NotImplemented("operation transfer.ListFilesHandlerFunc has not yet been implemented")
		})
	}

	redisClient := data.NewRedisCache(c.Config.Redis)
	dbClient := data.NewDatabaseClient(c.Config.Database)
	dataAccess := data.SyncheData{Cache: redisClient, Database: dbClient}

	api.TransferUploadChunkHandler = transfer.UploadChunkHandlerFunc(func(params transfer.UploadChunkParams) middleware.Responder {
		return handlers.UploadChunkHandler(params, dataAccess)
	})

	api.TransferNewUploadHandler = transfer.NewUploadHandlerFunc(func(params transfer.NewUploadParams) middleware.Responder {
		return handlers.NewUploadFileHandler(params, dataAccess)
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
