package main

import (
	"github.com/asaberwd/atomica-blog/handlers/docs"
	"github.com/asaberwd/atomica-blog/handlers/health"
	"github.com/asaberwd/atomica-blog/swagger/restapi"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations"
	"github.com/go-openapi/loads"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		panic(err)
	}

	api := operations.NewAtomicaBlogServiceAPI(swaggerSpec)

	healthService := health.New()

	/*	db, err := sqlx.Connect("postgres", os.Getenv("PGCONN"))
		if err != nil {
			panic(err)
		}*/
	health.Configure(api, healthService)
	docs.Configure(api)
	//api.Serve(nil)

	server := restapi.NewServer(api)
	server.EnabledListeners = []string{"http"}
	server.Port = 8040
	//defer server.Shutdown()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
