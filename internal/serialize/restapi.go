package serialize

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/josestg/justforfun/pkg/mux"
)

// RestAPI encodes the given data to JSON and write it the given w.
func RestAPI(ctx context.Context, w http.ResponseWriter, data interface{}, status int) error {
	// If the context is missing this value, this is a serious problem,
	// because Mux Handle is never executed.
	v, err := mux.GetState(ctx)
	if err != nil {
		return mux.NewShutdownError(err.Error())
	}

	// Add status code into ContextValue.
	// So, the next/after middleware can use it.
	v.StatusCode = status

	// If there is nothing to marshal then set status code and return.
	if status == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	// Encode the data to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(status)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
