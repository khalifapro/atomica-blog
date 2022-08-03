package main

import (
	"log"
	"os"

	comments "github.com/asaberwd/atomica-blog/handlers/comment"
	"github.com/asaberwd/atomica-blog/handlers/docs"
	"github.com/asaberwd/atomica-blog/handlers/health"
	posts "github.com/asaberwd/atomica-blog/handlers/post"
	commentService "github.com/asaberwd/atomica-blog/internal/comment"
	postService "github.com/asaberwd/atomica-blog/internal/post"
	"github.com/asaberwd/atomica-blog/swagger/restapi"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		panic(err)
	}

	api := operations.NewAtomicaBlogServiceAPI(swaggerSpec)

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

	server := restapi.NewServer(api)
	server.EnabledListeners = []string{"http", "https"}
	server.Port = 8080
	defer server.Shutdown()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
