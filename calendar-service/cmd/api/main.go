package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/config"
	"github.com/NextGenPlatformUnni/calendar-service/internal/application/command"
	"github.com/NextGenPlatformUnni/calendar-service/internal/application/query"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/service"
	infrahttp "github.com/NextGenPlatformUnni/calendar-service/internal/infrastructure/http"
	"github.com/NextGenPlatformUnni/calendar-service/internal/infrastructure/postgres"
	"github.com/NextGenPlatformUnni/calendar-service/internal/presentation/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// DB
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Infrastructure
	recordRepo := postgres.NewRecordRepository(pool)
	scheduleRepo := postgres.NewScheduleRepository(pool)
	cycleClient := infrahttp.NewCycleRuleClient(cfg.AdminAPIURL, 3*time.Second)
	notifyClient := infrahttp.NewMockNotificationClient()

	// Domain
	calculator := service.NewCycleCalculator()

	// Commands
	calcScheduleHandler := command.NewCalculateScheduleCommandHandler(scheduleRepo, cycleClient, calculator)
	createRecordHandler := command.NewCreateRecordCommandHandler(recordRepo, calcScheduleHandler)
	updateRecordHandler := command.NewUpdateRecordCommandHandler(recordRepo, scheduleRepo, calcScheduleHandler)
	deleteRecordHandler := command.NewDeleteRecordCommandHandler(recordRepo, scheduleRepo)
	reservationHandler := command.NewHandleReservationFixedCommandHandler(recordRepo, calcScheduleHandler)
	completeScheduleHandler := command.NewCompleteScheduleCommandHandler(scheduleRepo)
	deleteScheduleHandler := command.NewDeleteScheduleCommandHandler(scheduleRepo)
	_ = command.NewProcessRemindersCommandHandler(scheduleRepo, recordRepo, notifyClient)

	// Queries
	listRecordsHandler := query.NewListRecordsQueryHandler(recordRepo)
	getRecordHandler := query.NewGetRecordQueryHandler(recordRepo)
	listSchedulesHandler := query.NewListSchedulesQueryHandler(scheduleRepo)
	getScheduleHandler := query.NewGetScheduleQueryHandler(scheduleRepo)
	getStatsHandler := query.NewGetStatisticsQueryHandler(recordRepo)

	// Echo
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// Health
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok", "service": "calendar-service"})
	})

	// Routes
	api := e.Group("/api")

	recordHandler := handler.NewRecordHandler(createRecordHandler, updateRecordHandler, deleteRecordHandler, listRecordsHandler, getRecordHandler)
	recordHandler.Register(api.Group("/records"))

	scheduleHandler := handler.NewScheduleHandler(completeScheduleHandler, deleteScheduleHandler, listSchedulesHandler, getScheduleHandler)
	scheduleHandler.Register(api.Group("/schedules"))

	statsHandler := handler.NewStatisticsHandler(getStatsHandler)
	statsHandler.Register(api.Group("/statistics"))

	treatmentDataHandler := handler.NewTreatmentDataHandler(cfg.AdminAPIURL)
	treatmentDataHandler.Register(api.Group("/treatment-data"))

	mockHandler := handler.NewMockHandler(reservationHandler)
	mockHandler.Register(api.Group("/mock"))

	// Graceful shutdown
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
