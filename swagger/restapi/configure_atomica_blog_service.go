// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"

	comments "github.com/asaberwd/atomica-blog/handlers/comment"
	"github.com/asaberwd/atomica-blog/handlers/docs"
	"github.com/asaberwd/atomica-blog/handlers/health"
	posts "github.com/asaberwd/atomica-blog/handlers/post"
	commentService "github.com/asaberwd/atomica-blog/internal/comment"
	postService "github.com/asaberwd/atomica-blog/internal/post"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations/post"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:generate swagger generate server --target ../../swagger --name AtomicaBlogService --spec ../api.yaml --principal interface{}

func configureFlags(api *operations.AtomicaBlogServiceAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AtomicaBlogServiceAPI) http.Handler {
	// configure the api here
	swaggerSpec, err := loads.Analyzed(SwaggerJSON, "")
	if err != nil {
		log.Fatal("Invalid swagger file for initializing user", err)
	}
	api = operations.NewAtomicaBlogServiceAPI(swaggerSpec)
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.HTMLProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("html producer has not yet been implemented")
	})
	api.JSONProducer = runtime.JSONProducer()

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// post.UpdatePostWithFormMaxParseMemory = 32 << 20

	if api.PostUpdatePostHandler == nil {
		api.PostUpdatePostHandler = post.UpdatePostHandlerFunc(func(params post.UpdatePostParams) middleware.Responder {
			return middleware.NotImplemented("operation post.UpdatePost has not yet been implemented")
		})
	}

	healthService := health.New()
	health.Configure(api, healthService)

	docs.Configure(api)

	db, err := sqlx.Connect("postgres", os.Getenv("PGCONN"))
	if err != nil {
		panic(err)
	}
	postManager := postService.NewManager(db)
	postSvc := posts.New(postManager)
	posts.Configure(api, *postSvc)

	commentManager := commentService.NewManager(db)
	CommentSvc := comments.New(commentManager)
	comments.Configure(api, *CommentSvc)

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
