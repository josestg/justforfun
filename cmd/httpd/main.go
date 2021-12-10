package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/josestg/justforfun/internal/domain/sys"

	"github.com/josestg/justforfun/internal/delivery/restapi"
)

var (
	buildRef  = "unknown"
	buildDate = "unknown"
	buildName = "unknown"
)

func main() {
	sys.BuildRef.Set(buildRef)
	sys.BuildName.Set(buildName)
	sys.BuildDate.Set(buildDate)

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}

func run() error {
	logger := log.New(os.Stdout, "HTTPD ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	logger.Println("main:", "started")
	defer logger.Println("main:", "stopped")

	// create web server for our API.
	//
	// To notify the server to shutting down gracefully when the shutdownChannel
	// receives a termination signal or an interrupt signal.
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)

	router := restapi.NewRouter(&restapi.Option{
		Logger:          logger,
		ShutdownChannel: shutdownChannel,
	})

	server := &http.Server{
		Handler:      router,
		Addr:         ":3000",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	listenErr := make(chan error, 1)
	go func() {
		logger.Println("main:", "server is listening on "+server.Addr)
		listenErr <- server.ListenAndServe()
	}()

	select {
	case sig := <-shutdownChannel:
		logger.Println("main:", "server receives shutdown signal:", sig)

		// creating a deadline for the server to complete the incoming request
		// before the shutdown signal is received.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Gracefully shutdown
		// a non-nil error caused by the server failing to shut down gracefully
		// or the server exceeding the shutdown timeout.
		if err := server.Shutdown(ctx); err != nil {
			// Force shutdown
			return server.Close()
		}

		logger.Println("main:", "server was shutting down gracefully")
	case err := <-listenErr:
		if err != nil {
			return fmt.Errorf("%w: listening failed", err)
		}
	}

	return nil
}
