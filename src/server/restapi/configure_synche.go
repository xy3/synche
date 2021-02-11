// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/handlers"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/testing"
	"log"
	"net/http"
)

//go:generate swagger generate server --target ../../server --name Synche --spec ../api/openapi-spec/synche-server-api.yaml --principal models.Message

func configureFlags(api *operations.SyncheAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SyncheAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	err := c.InitConfig()
	if err != nil {
		log.Printf("Fatal error config file: %s \n", err)
	}

	// Read updates to the config file while server is running
	viper.WatchConfig()

	// Config vars
	dbDriver := viper.GetString("database.driver")
	dbUsername := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	dbProtocol := viper.GetString("database.protocol")
	dbAddress := viper.GetString("database.address")
	dbName := viper.GetString("database.name")

	// Create database with chunk table and connection_request table if they don't exist
	err = database.CreateDatabase(dbDriver, dbUsername, dbPassword, dbProtocol, dbAddress, dbName)
	if err != nil {
		log.Printf("Database creation failed: %s", err)
	}

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// files.UploadFileMaxParseMemory = 32 << 20

	// ============= Start Route Handlers =============
	api.TestingCheckGetHandler = testing.CheckGetHandlerFunc(handlers.CheckGetHandler)

	// TODO: Implement listing functionality
	if api.FilesListFilesHandler == nil {
		api.FilesListFilesHandler = files.ListFilesHandlerFunc(func(params files.ListFilesParams) middleware.Responder {
			return middleware.NotImplemented("operation files.ListFilesHandlerFunc has not yet been implemented")
		})
	}

	clientBuilder := database.NewDBClientBuilder()
	dbClient := clientBuilder.BuildSqlClient(dbDriver, dbUsername, dbPassword, dbProtocol, dbAddress, dbName)

	api.FilesUploadChunkHandler = files.UploadChunkHandlerFunc(func(params files.UploadChunkParams) middleware.Responder {
		return handlers.UploadChunkHandler(params, dbClient)})


	api.FilesNewUploadHandler = files.NewUploadHandlerFunc(func(params files.NewUploadParams) middleware.Responder {
		return handlers.NewUploadFileHandler(params, dbClient)})
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
