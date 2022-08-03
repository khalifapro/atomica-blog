package docs

import (
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations/doc"
	"github.com/go-openapi/runtime/middleware"
)

// Configure swagger docs API
func Configure(api *operations.AtomicaBlogServiceAPI) {
	api.DocGetDocHandler = doc.GetDocHandlerFunc(func(params doc.GetDocParams) middleware.Responder {
		return NewGetDocOK()
	})
}
