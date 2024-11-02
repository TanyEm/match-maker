package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/gin-gonic/gin"
)

type ServiceConfig struct {
	Port             int           `env:"PORT" envDefault:"8080"`
	ShutdownDuration time.Duration `env:"SHUTDOWN_DURATION" envDefault:"3s"`
}

func main() {
	cfg := ServiceConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	err := run(ctx, &cfg)
	cancel()

	if err != nil {
		log.Fatalf("Run returned an error: %v", err)
	}
}

func run(ctx context.Context, cfg *ServiceConfig) error {
	errCh := make(chan error)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		errStart := srv.ListenAndServe()

		if errors.Is(errStart, http.ErrServerClosed) {
			return
		}

		if errStart != nil {
			errCh <- fmt.Errorf("failed to start server: %w", errStart)
		}
	}()

	select {
	case err := <-errCh:
		log.Printf("Error starting server: %v", err)
	default:
		<-ctx.Done()
		log.Printf("Shutting down server...")
	}

	log.Printf("Shutting down server in %s", cfg.ShutdownDuration.String())

	time.Sleep(cfg.ShutdownDuration)

	srv.Shutdown(ctx)

	return nil
}
