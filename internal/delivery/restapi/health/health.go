package health

import (
	"fmt"
	"net/http"

	dHealth "github.com/josestg/justforfun/internal/domain/health"

	"github.com/josestg/justforfun/internal/serialize"
)

// Handler is a health handler.
// This handler serves APIs for checking system health.
type Handler struct {
	u dHealth.UseCase
}

// NewHandler creates a new health handler.
func NewHandler(u dHealth.UseCase) *Handler {
	return &Handler{
		u: u,
	}
}

func (h *Handler) showHealthStatus(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	report, err := h.u.HealthReport(ctx)
	if err != nil {
		return fmt.Errorf("%w: getting sys health repost", err)
	}

	return serialize.RestAPI(ctx, w, report, http.StatusOK)
}

// ServeHTTP serves the Health Handler at /v1/healths.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	default:
		http.NotFound(w, r)
		return nil
	case http.MethodGet:
		return h.showHealthStatus(w, r)
	}
}
