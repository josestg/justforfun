package restapi

import (
	"log"

	"github.com/josestg/justforfun/internal/delivery/restapi/health"

	"github.com/josestg/justforfun/pkg/mux"
)

// Option contains all required dependencies to serve the HTTP REST API delivery.
type Option struct {
	Logger          *log.Logger
	ShutdownChannel mux.ShutdownChannel
}

// NewRouter creates a configured router for HTTP REST API delivery.
func NewRouter(opt *Option) *mux.Router {
	router := mux.NewRouter(opt.ShutdownChannel)

	healthHandler := health.NewHandler()

	router.Handle("/v1/healths", healthHandler)

	return router
}
