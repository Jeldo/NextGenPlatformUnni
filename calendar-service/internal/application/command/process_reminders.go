package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
)

type ProcessRemindersCommand struct {
	Date time.Time
}

type ProcessRemindersCommandHandler struct {
	scheduleRepo   port.ScheduledTreatmentRepository
	recordRepo     port.TreatmentRecordRepository
	notifyClient   port.NotificationClient
}

func NewProcessRemindersCommandHandler(scheduleRepo port.ScheduledTreatmentRepository, recordRepo port.TreatmentRecordRepository, notifyClient port.NotificationClient) *ProcessRemindersCommandHandler {
	return &ProcessRemindersCommandHandler{scheduleRepo: scheduleRepo, recordRepo: recordRepo, notifyClient: notifyClient}
}

func (h *ProcessRemindersCommandHandler) Handle(ctx context.Context, cmd ProcessRemindersCommand) error {
	now := cmd.Date
	schedules, err := h.scheduleRepo.FindDueSchedules(ctx, now)
	if err != nil {
		return model.NewAppError(model.ErrInternal, "failed to find due schedules", err)
	}

	for _, s := range schedules {
		// Check if user already has a future record for this treatment
		futureRecords, _ := h.scheduleRepo.FindFutureByUserAndTreatment(ctx, s.UserID, s.TreatmentID, now)
		if len(futureRecords) > 0 {
			continue // skip reminder if already booked
		}

		daysSince := int(now.Sub(s.CreatedAt).Hours() / 24)
		msg := s.BuildReminderMessage(daysSince)

		if err := h.notifyClient.SendReminder(ctx, msg); err != nil {
			slog.Error("failed to send reminder", "schedule_id", s.ID, "error", err)
			continue
		}

		s.MarkReminded(now)
		if err := h.scheduleRepo.Update(ctx, s); err != nil {
			slog.Error("failed to update schedule status", "schedule_id", s.ID, "error", err)
		}
	}

	return nil
}
