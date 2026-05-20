package query

import (
	"context"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
)

type ListRecordsQuery struct {
	UserID string
	From   time.Time
	To     time.Time
}

type ListRecordsQueryHandler struct {
	repo port.TreatmentRecordRepository
}

func NewListRecordsQueryHandler(repo port.TreatmentRecordRepository) *ListRecordsQueryHandler {
	return &ListRecordsQueryHandler{repo: repo}
}

func (h *ListRecordsQueryHandler) Handle(ctx context.Context, q ListRecordsQuery) ([]*model.TreatmentRecord, error) {
	records, err := h.repo.FindByUserIDAndDateRange(ctx, q.UserID, q.From, q.To)
	if err != nil {
		return nil, model.NewAppError(model.ErrInternal, "failed to list records", err)
	}
	return records, nil
}
