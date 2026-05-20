package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/config"
	"github.com/NextGenPlatformUnni/calendar-service/internal/application/command"
	infrahttp "github.com/NextGenPlatformUnni/calendar-service/internal/infrastructure/http"
	"github.com/NextGenPlatformUnni/calendar-service/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	recordRepo := postgres.NewRecordRepository(pool)
	scheduleRepo := postgres.NewScheduleRepository(pool)
	notifyClient := infrahttp.NewMockNotificationClient()

	remindersHandler := command.NewProcessRemindersCommandHandler(scheduleRepo, recordRepo, notifyClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	slog.Info("cron scheduler started")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				slog.Info("processing reminders")
				if err := remindersHandler.Handle(ctx, command.ProcessRemindersCommand{Date: time.Now()}); err != nil {
					slog.Error("reminder processing failed", "error", err)
				}
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("cron scheduler shutting down")
}
