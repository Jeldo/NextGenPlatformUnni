package http

import (
	"context"
	"log/slog"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
)

type MockNotificationClient struct{}

func NewMockNotificationClient() *MockNotificationClient {
	return &MockNotificationClient{}
}

func (c *MockNotificationClient) SendReminder(ctx context.Context, msg model.ReminderMessage) error {
	slog.Info("sending reminder (mock)", "user_id", msg.UserID, "treatment", msg.TreatmentName, "days_since", msg.DaysSince)
	return nil
}
