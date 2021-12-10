package mux

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"time"
)

var (
	// ErrStateAreMissing is an error when the request state not exist in
	// request context.
	ErrStateAreMissing = errors.New("state are missing from request context")
)

// keyType represents the type of value for the context key
type keyType int

// StateKey is a key to stores and retrieves the State
// from the request context.
const StateKey = keyType(0)

// State is the initial state for each request.
type State struct {
	StatusCode     int
	RequestCreated time.Time
}

// GetState gets the initial state form the given context.
func GetState(ctx context.Context) (*State, error) {
	v, valid := ctx.Value(StateKey).(*State)
	if !valid {
		return nil, ErrStateAreMissing
	}

	return v, nil
}

// Handler is just a http.Handler that can returns an error.
// By returning an error, now we can make a centralized error handling.
type Handler interface {
	// ServeHTTP is http.Handler but returns an error.
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP calls h(w,r).
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

// Middleware is a function that will be executed before or/and after the given
// handler has been executed.
//
// There are two types of middleware in this Router.
// The Global middlewares and The Route middlewares.
// The Global middlewares will be applied to all handler, meanwhile the Route
// middlewares only applied specific to the handler in the given route.
type Middleware func(handler Handler) Handler

// applyMiddleware applies the middleware to handler by wrapping the origin
// handler with the given middlewares.
func applyMiddleware(handler Handler, middlewares []Middleware) Handler {
	// Using a backward loop is to make the first element in the slice to be
	// executed first. This wrapping process expect the original handler
	// to be a center of the layers and the first middleware to be an outer layer.
	//
	// Visually, the final layer will look like this.
	//
	// M[0](
	//   M[1](
	//     ...
	//       M[k-2](
	//         M[k-1](handler)
	//       )
	//      ...
	//    )
	// )

	// If M is a middleware, M[i] is a middleware at index i in middlewares slice.
	// Where i is in range [k-1, 0] and k is length of slice.
	for i := len(middlewares) - 1; i >= 0; i-- {
		fn := middlewares[i]
		if fn != nil {
			handler = fn(handler)
		}
	}
	return handler
}

type Router struct {
	mux *http.ServeMux
	sc  ShutdownChannel
	gm  []Middleware
}

func NewRouter(channel ShutdownChannel, middleware ...Middleware) *Router {
	mux := http.NewServeMux()
	return &Router{
		mux: mux,
		sc:  channel,
		gm:  middleware,
	}
}

// SignalShutdown sends a shutdown signal through the shutdown channel.
func (r *Router) SignalShutdown() {
	if r.sc != nil {
		r.sc <- syscall.SIGTERM
	}
}

func (r *Router) Handle(pattern string, handler Handler, middleware ...Middleware) {
	// wraps original handler with given middlewares.
	handler = applyMiddleware(handler, middleware)
	// wraps the wrapped original handler again with r.gm.
	handler = applyMiddleware(handler, r.gm)

	fn := func(w http.ResponseWriter, req *http.Request) {
		// create the initial state
		s := State{
			StatusCode:     http.StatusOK, // default status code
			RequestCreated: time.Now(),
		}

		// Set an initial value for each request.
		ctx := context.WithValue(req.Context(), StateKey, &s)
		if err := handler.ServeHTTP(w, req.WithContext(ctx)); err != nil {
			// makes a shutdown signal if a critical error occurred.
			if IsShutdownError(err) {
				r.SignalShutdown()
			}
		}
	}

	r.mux.Handle(pattern, http.HandlerFunc(fn))
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(rw, req)
}
