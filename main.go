package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	loginapi "awesome/login/api"
	smurfsapi "awesome/smurfs/api"
)

const (
	serverReadHeaderTimeout = 5 * time.Second
	serverWriteTimeout      = 10 * time.Second
	serverReadTimeOut       = 120 * time.Second
)

type application struct {
	*http.Server
}

func newApplication() *application {
	handler := mountModules()
	return &application{Server: &http.Server{
		Addr:              ":8080",
		ReadTimeout:       serverReadTimeOut,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		WriteTimeout:      serverWriteTimeout,
		Handler:           handler,
	}}
}

func mountModules() chi.Router {
	router := coreRouter()
	if err := smurfsapi.Mount(router); err != nil {
		log.Panic().Err(err).Msg("Unable to mount smurfs module.")
	}
	if err := loginapi.Mount(router); err != nil {
		log.Panic().Err(err).Msg("Unable to mount login module.")
	}

	dumpRoutes(router)
	return router
}

func dumpRoutes(router chi.Router) {
	log.Debug().Msg("Dumping chi routes:")
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Debug().Msgf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Debug().Msgf("Logging err: %s\n", err.Error())
	}
}

func coreRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))
	router.Use(middleware.RealIP)
	router.Use(middleware.GetHead)
	router.Use(middleware.Timeout(2 * time.Second))

	router.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
	router.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	router.Use(middleware.SetHeader("X-DNS-Prefetch-Control", "off"))
	router.Use(middleware.SetHeader("Referrer-Policy", "no-referrer"))

	return router
}

func (a *application) start() {
	log.Logger.Info().Msg("Application starting.")

	a.startHTTPServer()

	log.Logger.Info().Msg("Application started.")
}

func (a *application) startHTTPServer() {
	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Logger.Info().Msg("Server shutdown.")
			} else {
				log.Logger.Panic().Err(err).Msg("Unexpected server error.")
			}
		}
	}()
}

func (a *application) stop() error {
	log.Logger.Info().Msg("Application shutting down.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := a.Server.Shutdown(ctx)

	log.Logger.Info().Msg("Application shutdown.")
	return err
}

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	_ = os.Setenv("TZ", "UTC")
	setupLogging()

	app := newApplication()
	app.start()
	<-shutdown
	_ = app.stop()
}

func setupLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.With().Caller().Str("application", "awesome").Logger()
}
