package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func (app *Application) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		log.Info().Str("signal", s.String()).Msg("shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Info().Str("addr", srv.Addr).Msg("completing background tasks")

		app.wg.Wait()
		shutdownError <- nil
	}()

	log.Info().Str("addr", srv.Addr).Msg("starting server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Info().Str("addr", srv.Addr).Msg("stopped server")

	return nil
}
