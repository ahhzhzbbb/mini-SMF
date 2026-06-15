package main

import (
	"context"
	"fmt"
	"io"
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/proxy"
	"mini-SMF/gateway/internal/registry"
	"mini-SMF/gateway/internal/router"
	"mini-SMF/pkg/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	//config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	//logger
	logger := logger.NewLogger(cfg.LogLevel)

	//registry
	registry := registry.NewRegistry()
	registry.Load(os.Getenv("PDU_SERVICE_NAME"))

	//loadbalancer (Round Robin)
	var loadBalancer router.LoadBalancer = router.NewRoundRobin(0)

	proxy := proxy.NewProxy(
		cfg,
		logger,
		registry,
		loadBalancer,
	)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: proxy,
	}

	go func() {
		fmt.Fprintf(w, "listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)

		}
	}()

	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Println("your gateway is shutting down...")
		os.Exit(1)
	}

}
