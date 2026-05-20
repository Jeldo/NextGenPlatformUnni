package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/application/command"
	"github.com/NextGenPlatformUnni/calendar-service/internal/application/query"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/presentation/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RecordHandler struct {
	createHandler *command.CreateRecordCommandHandler
	updateHandler *command.UpdateRecordCommandHandler
	deleteHandler *command.DeleteRecordCommandHandler
	listHandler   *query.ListRecordsQueryHandler
	getHandler    *query.GetRecordQueryHandler
}

func NewRecordHandler(
	createHandler *command.CreateRecordCommandHandler,
	updateHandler *command.UpdateRecordCommandHandler,
	deleteHandler *command.DeleteRecordCommandHandler,
	listHandler *query.ListRecordsQueryHandler,
	getHandler *query.GetRecordQueryHandler,
) *RecordHandler {
	return &RecordHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		listHandler:   listHandler,
		getHandler:    getHandler,
	}
}

func (h *RecordHandler) Register(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *RecordHandler) Create(c echo.Context) error {
	var req dto.CreateRecordRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
	}

	treatmentID, _ := uuid.Parse(req.TreatmentID)
	categoryID, _ := uuid.Parse(req.CategoryID)
	treatmentDate, _ := time.Parse(time.RFC3339, req.TreatmentDate)

	input := command.CreateRecordCommand{
		UserID:        req.UserID,
		TreatmentID:   treatmentID,
		TreatmentName: req.TreatmentName,
		CategoryID:    categoryID,
		CategoryName:  req.CategoryName,
		TreatmentDate: treatmentDate,
		HospitalName:  req.HospitalName,
	}

	if req.DosageValue != nil {
		v, _ := strconv.ParseFloat(*req.DosageValue, 64)
		input.DosageValue = &v
	}
	if req.DosageUnit != nil {
		u := model.DosageUnit(*req.DosageUnit)
		input.DosageUnit = &u
	}

	record, err := h.createHandler.Handle(c.Request().Context(), input)
	if err != nil {
		return handleAppError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.ToRecordResponse(record))
}

func (h *RecordHandler) List(c echo.Context) error {
	userID := c.QueryParam("user_id")
	from, _ := time.Parse(time.RFC3339, c.QueryParam("from"))
	to, _ := time.Parse(time.RFC3339, c.QueryParam("to"))

	if to.IsZero() {
		to = time.Now().AddDate(1, 0, 0)
	}

	records, err := h.listHandler.Handle(c.Request().Context(), query.ListRecordsQuery{
		UserID: userID, From: from, To: to,
	})
	if err != nil {
		return handleAppError(c, err)
	}

	resp := make([]dto.RecordResponse, len(records))
	for i, r := range records {
		resp[i] = dto.ToRecordResponse(r)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *RecordHandler) Get(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))
	userID := c.QueryParam("user_id")

	record, err := h.getHandler.Handle(c.Request().Context(), query.GetRecordQuery{ID: id, UserID: userID})
	if err != nil {
		return handleAppError(c, err)
	}

	return c.JSON(http.StatusOK, dto.ToRecordResponse(record))
}

func (h *RecordHandler) Update(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))
	var req dto.UpdateRecordRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
	}

	treatmentID, _ := uuid.Parse(req.TreatmentID)
	categoryID, _ := uuid.Parse(req.CategoryID)
	treatmentDate, _ := time.Parse(time.RFC3339, req.TreatmentDate)

	input := command.UpdateRecordCommand{
		ID:            id,
		UserID:        c.QueryParam("user_id"),
		TreatmentID:   treatmentID,
		TreatmentName: req.TreatmentName,
		CategoryID:    categoryID,
		CategoryName:  req.CategoryName,
		TreatmentDate: treatmentDate,
		HospitalName:  req.HospitalName,
	}

	if req.DosageValue != nil {
		v, _ := strconv.ParseFloat(*req.DosageValue, 64)
		input.DosageValue = &v
	}
	if req.DosageUnit != nil {
		u := model.DosageUnit(*req.DosageUnit)
		input.DosageUnit = &u
	}

	record, err := h.updateHandler.Handle(c.Request().Context(), input)
	if err != nil {
		return handleAppError(c, err)
	}

	return c.JSON(http.StatusOK, dto.ToRecordResponse(record))
}

func (h *RecordHandler) Delete(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))
	userID := c.QueryParam("user_id")

	if err := h.deleteHandler.Handle(c.Request().Context(), command.DeleteRecordCommand{ID: id, UserID: userID}); err != nil {
		return handleAppError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func handleAppError(c echo.Context, err error) error {
	var appErr *model.AppError
	if errors.As(err, &appErr) {
		status := http.StatusInternalServerError
		switch appErr.Code {
		case model.ErrValidation:
			status = http.StatusBadRequest
		case model.ErrNotFound:
			status = http.StatusNotFound
		case model.ErrConflict:
			status = http.StatusConflict
		case model.ErrForbidden:
			status = http.StatusForbidden
		case model.ErrExternalAPI:
			status = http.StatusBadGateway
		}
		return respondError(c, status, string(appErr.Code), appErr.Message)
	}
	return respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected error")
}

func respondError(c echo.Context, status int, code, message string) error {
	return c.JSON(status, dto.ErrorResponse{
		Error: dto.ErrorDetail{Code: code, Message: message},
	})
}
