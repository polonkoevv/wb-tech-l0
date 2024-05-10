package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/stan.go"
	"github.com/polonkoevv/wb-tech/internal/api"
	"github.com/polonkoevv/wb-tech/internal/config"
	"github.com/polonkoevv/wb-tech/internal/pkg/logger"
	"github.com/polonkoevv/wb-tech/internal/service"
	"github.com/polonkoevv/wb-tech/internal/storage"
	"github.com/polonkoevv/wb-tech/internal/storage/postgres"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Level)

	log.Info("config", slog.Any("config", cfg))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc, err := stan.Connect(cfg.Nats.ClusterID, cfg.Nats.ClientId, stan.NatsURL(cfg.Nats.ListenUrl))
	if err != nil {
		log.Error("Failed to connect to NATS Streaming: %v", err)
		os.Exit(1)
	}

	conn, err := storage.OpenDB(ctx, cfg.Storage)

	db := postgres.New(conn)

	if err != nil {
		log.Error("Opening DB connection :", slog.Any("error", err.Error()))
	}

	srv := service.New(sc, db, log)

	err = srv.LoadCache(ctx)
	if err != nil {
		log.Error(err.Error())
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return srv.Listen(ctx, cfg.Nats.ListenChanel)
	})

	r := api.New(srv)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPServer.Port,
		Handler: r,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed to listen: %v", err)
			os.Exit(1)
		}
	}()

	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Warn("exit reason: %s \n", err)
	}

}
