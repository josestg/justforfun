package mux

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"syscall"
	"testing"
	"time"
)

func TestGetState(t *testing.T) {
	if _, err := GetState(context.Background()); err == nil {
		t.Errorf("expecting error nil but got: %v", err)
	}

	ctx := context.WithValue(context.Background(), StateKey, &State{
		RequestCreated: time.Now(),
		StatusCode:     0,
	})

	// Try to wrap context with another context.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if _, err := GetState(ctx); err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}
}

func TestRouter_Handle(t *testing.T) {
	const (
		expectedResponseBody = "TEST_OK"
		expectedRequestID    = "123"
		exampleURL           = "/example"
	)

	t.Run("expecting ok", func(t *testing.T) {
		router := NewRouter(nil)

		handler := func(w http.ResponseWriter, r *http.Request) error {
			_, err := GetState(r.Context())
			if err != nil {
				t.Fatal("expecting request has initial state")
			}

			_, err = io.WriteString(w, expectedResponseBody)
			return err
		}

		// register handler
		router.Handle(exampleURL, HandlerFunc(handler))

		// make a http request
		req := httptest.NewRequest(http.MethodPost, exampleURL, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		statusCode := rec.Result().StatusCode
		if statusCode != http.StatusOK {
			t.Errorf("expecting status code: %d but got: %d", http.StatusOK, statusCode)
		}

		body := rec.Body.String()
		if body != expectedResponseBody {
			t.Errorf("expecting response body: %s but got: %s", expectedRequestID, body)
		}
	})

	t.Run("expecting got a shutdown signal", func(t *testing.T) {
		shutdownChannel := make(ShutdownChannel, 1)
		router := NewRouter(shutdownChannel)

		handler := func(w http.ResponseWriter, r *http.Request) error {
			return NewShutdownError("fake shutdown error")
		}

		// register handler
		router.Handle(exampleURL, HandlerFunc(handler))

		// make a http request
		req := httptest.NewRequest(http.MethodPost, exampleURL, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		select {
		case sig := <-shutdownChannel:
			if sig != syscall.SIGTERM {
				t.Fatalf("expecting termination signal but got: %v", sig)
			}
		// If no signal received in next one second.
		case <-time.After(time.Second):
			t.Fatalf("expecting got a shutdwon signal")
		}
	})
}

func TestRouter_MiddlewareHandle(t *testing.T) {
	// Objective:
	// Check the execution order.
	//
	// Expect:
	// Global Middlewares = {gm0, gm1}
	// Router Middlewares = {rm0, rm1}
	// HandlerFunc = fn
	// The Call Stack should be look like this:
	// gm0: start
	// 	gm1: start
	// 		rm0: start
	//			rm1: start
	// 				fn: start
	//				fn: end
	// 			rm1: end
	//		rm0: end
	//	gm1: end
	// gm0:end
	callStack := make([]int, 0)

	// Middleware factory
	factory := func(flag int) Middleware {
		return func(handler Handler) Handler {
			fn := func(w http.ResponseWriter, r *http.Request) error {
				callStack = append(callStack, flag)
				defer func() {
					callStack = append(callStack, flag)
				}()

				return handler.ServeHTTP(w, r)
			}
			return HandlerFunc(fn)
		}
	}

	gm0 := factory(0)
	gm1 := factory(1)
	rm0 := factory(2)
	rm1 := factory(3)

	router := NewRouter(nil, gm0, gm1)

	fn := func(w http.ResponseWriter, r *http.Request) error {
		callStack = append(callStack, 4)
		return nil
	}

	const exampleURL = "/example"
	router.Handle(exampleURL, HandlerFunc(fn), rm0, rm1)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, exampleURL, nil)

	router.ServeHTTP(rec, req)

	expected := []int{0, 1, 2, 3, 4, 3, 2, 1, 0}
	if !reflect.DeepEqual(expected, callStack) {
		t.Fatalf("expecting: %v got: %v", expected, callStack)
	}
}
