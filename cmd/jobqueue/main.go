package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"github.com/tysenmeroniuk/jobqueue/internal/api"
	"github.com/tysenmeroniuk/jobqueue/internal/config"
	"github.com/tysenmeroniuk/jobqueue/internal/db"
	"github.com/tysenmeroniuk/jobqueue/internal/handler"
	"github.com/tysenmeroniuk/jobqueue/internal/queue"
	"github.com/tysenmeroniuk/jobqueue/internal/worker"
)

func main() {
	root := &cobra.Command{Use: "jobqueue", Short: "Distributed job queue"}
	root.AddCommand(apiCmd(), workerCmd(), migrateCmd())
	if err := root.Execute(); err != nil {
		slog.Error("command failed", "err", err)
		os.Exit(1)
	}
}

func apiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Run the HTTP API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			ctx, cancel := signalContext()
			defer cancel()

			pool, err := db.Connect(ctx, cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("db connect: %w", err)
			}
			defer pool.Close()

			q := queue.New(pool)
			srv := api.New(q)

			httpSrv := &http.Server{
				Addr:         net.JoinHostPort("", strconv.Itoa(cfg.Port)),
				Handler:      srv.Router(),
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 30 * time.Second,
			}

			go func() {
				<-ctx.Done()
				shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer shutdownCancel()
				_ = httpSrv.Shutdown(shutdownCtx)
			}()

			slog.Info("api listening", "port", cfg.Port)
			if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return err
			}
			return nil
		},
	}
}

func workerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Run a worker process",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			ctx, cancel := signalContext()
			defer cancel()

			pool, err := db.Connect(ctx, cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("db connect: %w", err)
			}
			defer pool.Close()

			q := queue.New(pool)
			reg := handler.NewRegistry()

			// TODO: register job handlers here, e.g.:
			//   reg.Register("send_email", sendEmailHandler)

			workerID := cfg.WorkerID
			if workerID == "" {
				hn, _ := os.Hostname()
				workerID = fmt.Sprintf("%s-%d", hn, os.Getpid())
			}

			w := worker.New(q, reg, workerID)
			r := worker.NewReaper(q, 30*time.Second)

			go func() {
				if err := w.Run(ctx); err != nil && err != context.Canceled {
					slog.Error("worker stopped", "err", err)
				}
			}()
			go func() {
				if err := r.Run(ctx); err != nil && err != context.Canceled {
					slog.Error("reaper stopped", "err", err)
				}
			}()

			slog.Info("worker running", "id", workerID)
			<-ctx.Done()
			return nil
		},
	}
}

func migrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate [up|down|status]",
		Short: "Run database migrations (goose)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			direction := "up"
			if len(args) > 0 {
				direction = args[0]
			}
			conn, err := sql.Open("pgx", cfg.DatabaseURL)
			if err != nil {
				return err
			}
			defer conn.Close()

			if err := goose.SetDialect("postgres"); err != nil {
				return err
			}
			return goose.Run(direction, conn, "migrations")
		},
	}
}

func signalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}
