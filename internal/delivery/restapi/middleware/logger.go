package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/josestg/justforfun/pkg/mux"
)

// Logger logs request information for each incoming request.
func Logger(logger *log.Logger) mux.Middleware {
	return func(handler mux.Handler) mux.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			state, err := mux.GetState(ctx)
			if err != nil {
				return mux.NewShutdownError(err.Error())
			}

			logger.Printf("logger: receiving: %s %s", r.Method, r.URL.Path)
			defer func(s *mux.State) {
				logger.Printf(
					"logger: completed: %s %s  %d  %s μs",
					r.Method, r.URL.Path, state.StatusCode, time.Since(state.RequestCreated),
				)
			}(state)

			return handler.ServeHTTP(w, r.WithContext(ctx))
		}

		return mux.HandlerFunc(fn)
	}
}
