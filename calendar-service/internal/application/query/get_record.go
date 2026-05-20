package query

import (
	"context"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type GetRecordQuery struct {
	ID     uuid.UUID
	UserID string
}

type GetRecordQueryHandler struct {
	repo port.TreatmentRecordRepository
}

func NewGetRecordQueryHandler(repo port.TreatmentRecordRepository) *GetRecordQueryHandler {
	return &GetRecordQueryHandler{repo: repo}
}

func (h *GetRecordQueryHandler) Handle(ctx context.Context, q GetRecordQuery) (*model.TreatmentRecord, error) {
	record, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		return nil, model.NewAppError(model.ErrNotFound, "record not found", err)
	}
	if !record.IsOwnedBy(q.UserID) {
		return nil, model.NewAppError(model.ErrForbidden, "not authorized", nil)
	}
	return record, nil
}
