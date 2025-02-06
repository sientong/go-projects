package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		oscall := <-ch
		log.Info().Msgf("Received signal: %s", oscall)
		cancel()
	}()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: newHttpHandler(),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Failed to start server")
		}
	}()

	<-ctx.Done()

	gracefulShutdownPeriod := time.Second * 10
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownPeriod)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown server gracefully")
	}

	log.Info().Msg("Server stopped")
}

func newHttpHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/binary", BinaryUpload())
	router.HandleFunc("/api/v1/form", FormUpload())
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})
	return router
}
