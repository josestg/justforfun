package health

import (
	"encoding/json"
	"net/http"
)

// Handler is a health handler.
// This handler serves APIs for checking system health.
type Handler struct {
}

// NewHandler creates a new health handler.
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) showHealthStatus(w http.ResponseWriter, _ *http.Request) error {
	data := struct {
		Ready bool `json:"ready"`
	}{
		Ready: true,
	}

	// Encode the data to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(http.StatusOK)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
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
