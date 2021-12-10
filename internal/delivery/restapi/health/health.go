package health

import (
	"net/http"

	"github.com/josestg/justforfun/internal/serialize"
)

// Handler is a health handler.
// This handler serves APIs for checking system health.
type Handler struct {
}

// NewHandler creates a new health handler.
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) showHealthStatus(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	data := struct {
		Ready bool `json:"ready"`
	}{
		Ready: true,
	}

	return serialize.RestAPI(ctx, w, &data, http.StatusOK)
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
