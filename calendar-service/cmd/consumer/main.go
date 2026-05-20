package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/config"
	"github.com/NextGenPlatformUnni/calendar-service/internal/application/command"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/service"
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
	cycleClient := infrahttp.NewCycleRuleClient(cfg.AdminAPIURL, 3*time.Second)
	calculator := service.NewCycleCalculator()

	calcHandler := command.NewCalculateScheduleCommandHandler(scheduleRepo, cycleClient, calculator)
	reservationHandler := command.NewHandleReservationFixedCommandHandler(recordRepo, calcHandler)

	// In-memory channel simulating SQS
	messages := make(chan []byte, 100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		slog.Info("consumer started, waiting for messages...")
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-messages:
				var event model.ReservationFixedEvent
				if err := json.Unmarshal(msg, &event); err != nil {
					slog.Error("failed to parse event", "error", err)
					continue
				}
				if err := reservationHandler.Handle(ctx, command.HandleReservationFixedCommand{Event: event}); err != nil {
					slog.Error("failed to handle event", "error", err)
				}
			}
		}
	}()

	_ = messages // expose for testing

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("consumer shutting down")
}
