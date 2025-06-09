package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultIdleTimeout    = 2 * time.Minute
	defaultReadTimeout    = 10 * time.Second
	defaultWriteTimeout   = 5 * time.Minute
	defaultShutdownPeriod = 30 * time.Second
)

// TODO: Get and Handle Upper context
func (glob *CommonGlobals) ServeHTTP(router *chi.Mux, port int) error {

	stdLogger, err := zap.NewStdLogAt(glob.Logger, zapcore.WarnLevel)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ErrorLog:     stdLogger,
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	glob.Logger.Info("starting server", zap.String("serveraddr", srv.Addr))
	// Glob.Logger.Info("starting server", slog.Group("server", "addr", srv.Addr))

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	glob.Logger.Info("starting server", zap.String("serveraddr", srv.Addr))

	glob.wg.Wait()
	return nil
}
