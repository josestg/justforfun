package restapi

import (
	"log"

	"github.com/josestg/justforfun/internal/usecase"

	"github.com/josestg/justforfun/internal/delivery/restapi/middleware"

	"github.com/josestg/justforfun/pkg/mux"
)

// Option contains all required dependencies to serve the HTTP REST API delivery.
type Option struct {
	Logger          *log.Logger
	UseCase         *usecase.Container
	ShutdownChannel mux.ShutdownChannel
}

// NewRouter creates a configured router for HTTP REST API delivery.
func NewRouter(opt *Option) *mux.Router {
	router := mux.NewRouter(
		opt.ShutdownChannel,
		middleware.Logger(opt.Logger),
		middleware.Panics(opt.Logger),
	)

	return router
}
