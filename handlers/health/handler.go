package health

import (
	log "github.com/asaberwd/atomica-blog/logging"

	"github.com/asaberwd/atomica-blog/swagger/restapi/operations"
	"github.com/asaberwd/atomica-blog/swagger/restapi/operations/health_api"
	"github.com/go-openapi/runtime/middleware"
)

// Configure ...
func Configure(api *operations.AtomicaBlogServiceAPI, service Service) {
	api.HealthAPIHealthCheckHandler = health_api.HealthCheckHandlerFunc(func(params health_api.HealthCheckParams) middleware.Responder {
		var nilRequestID *string
		requestID := log.GetRequestID(nilRequestID)
		service.SetServiceRequestID(requestID)
		result, err := service.GetHealth(params.HTTPRequest.Context())
		if err != nil {
			return health_api.NewHealthCheckServiceUnavailable()
		}
		return health_api.NewHealthCheckOK().WithXREQUESTID(requestID).WithPayload(result)
	})
}
