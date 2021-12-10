package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/josestg/justforfun/pkg/mux"
)

// Panics is middleware for panics recovery.
// This middleware transform panic into normal error.
func Panics(logger *log.Logger) mux.Middleware {
	return func(handler mux.Handler) mux.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) (err error) {
			ctx := r.Context()
			state, err := mux.GetState(ctx)
			if err != nil {
				return mux.NewShutdownError(err.Error())
			}

			defer func(state *mux.State) {
				if rec := recover(); rec != nil {
					err = fmt.Errorf("panics: %v", rec)

					logger.Printf(
						"panics: recovered: %s %s  %d  %s μs",
						r.Method, r.URL.Path, state.StatusCode, time.Since(state.RequestCreated),
					)
				}
			}(state)

			return handler.ServeHTTP(w, r.WithContext(ctx))
		}

		return mux.HandlerFunc(fn)
	}
}
