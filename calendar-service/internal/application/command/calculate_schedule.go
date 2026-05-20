package command

import (
	"context"
	"log/slog"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/service"
)

type CalculateScheduleCommand struct {
	Record *model.TreatmentRecord
}

type CalculateScheduleCommandHandler struct {
	scheduleRepo port.ScheduledTreatmentRepository
	cycleClient  port.CycleRuleClient
	calculator   *service.CycleCalculator
}

func NewCalculateScheduleCommandHandler(scheduleRepo port.ScheduledTreatmentRepository, cycleClient port.CycleRuleClient, calculator *service.CycleCalculator) *CalculateScheduleCommandHandler {
	return &CalculateScheduleCommandHandler{scheduleRepo: scheduleRepo, cycleClient: cycleClient, calculator: calculator}
}

func (h *CalculateScheduleCommandHandler) Handle(ctx context.Context, cmd CalculateScheduleCommand) (*model.ScheduledTreatment, error) {
	record := cmd.Record
	rule, err := h.cycleClient.GetByTreatmentID(ctx, record.TreatmentID)
	if err != nil {
		slog.Warn("failed to get cycle rule", "treatment_id", record.TreatmentID, "error", err)
		return nil, nil // graceful degradation
	}
	if rule == nil || rule.CycleDays <= 0 {
		return nil, nil
	}

	schedule := h.calculator.Calculate(record, rule.CycleDays)

	if err := h.scheduleRepo.Save(ctx, schedule); err != nil {
		return nil, model.NewAppError(model.ErrInternal, "failed to save schedule", err)
	}

	return schedule, nil
}
