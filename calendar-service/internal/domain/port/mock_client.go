package port

import (
	"context"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/google/uuid"
)

// MockCycleRuleClient implements CycleRuleClient for testing.
type MockCycleRuleClient struct {
	Rules map[uuid.UUID]*CycleRule
	Err   error
}

func NewMockCycleRuleClient() *MockCycleRuleClient {
	return &MockCycleRuleClient{Rules: make(map[uuid.UUID]*CycleRule)}
}

func (m *MockCycleRuleClient) GetByTreatmentID(ctx context.Context, treatmentID uuid.UUID) (*CycleRule, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	rule, ok := m.Rules[treatmentID]
	if !ok {
		return nil, nil
	}
	return rule, nil
}

// MockNotificationClient implements NotificationClient for testing.
type MockNotificationClient struct {
	Sent []model.ReminderMessage
	Err  error
}

func NewMockNotificationClient() *MockNotificationClient {
	return &MockNotificationClient{}
}

func (m *MockNotificationClient) SendReminder(ctx context.Context, msg model.ReminderMessage) error {
	if m.Err != nil {
		return m.Err
	}
	m.Sent = append(m.Sent, msg)
	return nil
}
