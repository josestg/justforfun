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

	"github.com/josestg/justforfun/internal/iam"
	"github.com/josestg/justforfun/internal/repository"

	"github.com/josestg/justforfun/internal/validation"
	"github.com/josestg/justforfun/internal/wording"
	"github.com/josestg/justforfun/internal/wording/locale"

	"github.com/josestg/justforfun/pkg/x"

	"github.com/josestg/justforfun/internal/usecase/provider"

	"github.com/josestg/justforfun/internal/usecase"

	"github.com/josestg/justforfun/internal/conf"

	"github.com/josestg/justforfun/pkg/pqx"

	"github.com/josestg/justforfun/pkg/xerrs"

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

	cfg := conf.New(
		conf.WithRestAPIFromOSEnv(),
		conf.WithDBPostgreFromOSEnv(),
	)

	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

}

func run(c *conf.Config) error {
	logger := log.New(os.Stdout, "HTTPD ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	logger.Println("main:", "started")
	defer logger.Println("main:", "stopped")

	// open database connection.
	//
	// note: we only open connection once,
	// if we need database connection we must pass it as dependency.
	db, err := pqx.Open(c.DB.Postgre)

	if err != nil {
		return xerrs.Wrap(err, "open database connection")
	}

	logger.Println("main:", "checking database connection")
	checkCtx, checkCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer checkCancel()

	if err := pqx.CheckConnection(checkCtx, 5, db); err != nil {
		return xerrs.Wrap(err, "checking database connection")
	}

	// setup UseCase.
	//
	// creates containerized use case.
	up := provider.Provider{
		Identifier: x.NewIdentifier(),
		Clock:      x.NewLocalClock(time.Local),
		Tokenizer:  iam.NewJwtProvider(iam.NewJwtRS256(nil)),
		Repository: repository.NewSQLContainer(db),
		Validator:  validation.NewValidator(wording.NewWording(locale.Dictionary)),
	}

	uc := usecase.NewContainer(&up)

	// create web server for our API.
	//
	// To notify the server to shutting down gracefully when the shutdownChannel
	// receives a termination signal or an interrupt signal.
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)

	router := restapi.NewRouter(&restapi.Option{
		Logger:          logger,
		ShutdownChannel: shutdownChannel,
		UseCase:         uc,
	})

	server := &http.Server{
		Handler:      router,
		Addr:         c.RestAPI.Addr,
		ReadTimeout:  c.RestAPI.ReadTimeout,
		WriteTimeout: c.RestAPI.WriteTimeout,
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
		ctx, cancel := context.WithTimeout(context.Background(), c.RestAPI.ShutdownTimeout)
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
			return xerrs.Wrap(err, "listening failed")
		}
	}

	return nil
}
