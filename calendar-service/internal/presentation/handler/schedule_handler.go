package handler

import (
	"net/http"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/application/command"
	"github.com/NextGenPlatformUnni/calendar-service/internal/application/query"
	"github.com/NextGenPlatformUnni/calendar-service/internal/presentation/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ScheduleHandler struct {
	completeHandler *command.CompleteScheduleCommandHandler
	deleteHandler   *command.DeleteScheduleCommandHandler
	listHandler     *query.ListSchedulesQueryHandler
	getHandler      *query.GetScheduleQueryHandler
}

func NewScheduleHandler(
	completeHandler *command.CompleteScheduleCommandHandler,
	deleteHandler *command.DeleteScheduleCommandHandler,
	listHandler *query.ListSchedulesQueryHandler,
	getHandler *query.GetScheduleQueryHandler,
) *ScheduleHandler {
	return &ScheduleHandler{
		completeHandler: completeHandler,
		deleteHandler:   deleteHandler,
		listHandler:     listHandler,
		getHandler:      getHandler,
	}
}

func (h *ScheduleHandler) Register(g *echo.Group) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.PATCH("/:id/complete", h.Complete)
	g.DELETE("/:id", h.Delete)
}

func (h *ScheduleHandler) List(c echo.Context) error {
	userID := c.QueryParam("user_id")
	from, _ := time.Parse(time.RFC3339, c.QueryParam("from"))
	to, _ := time.Parse(time.RFC3339, c.QueryParam("to"))

	if to.IsZero() {
		to = time.Now().AddDate(1, 0, 0)
	}

	schedules, err := h.listHandler.Handle(c.Request().Context(), query.ListSchedulesQuery{
		UserID: userID, From: from, To: to,
	})
	if err != nil {
		return handleAppError(c, err)
	}

	resp := make([]dto.ScheduleResponse, len(schedules))
	for i, s := range schedules {
		resp[i] = dto.ToScheduleResponse(s)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *ScheduleHandler) Get(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))
	userID := c.QueryParam("user_id")

	schedule, err := h.getHandler.Handle(c.Request().Context(), query.GetScheduleQuery{ID: id, UserID: userID})
	if err != nil {
		return handleAppError(c, err)
	}

	return c.JSON(http.StatusOK, dto.ToScheduleResponse(schedule))
}

func (h *ScheduleHandler) Complete(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))
	userID := c.QueryParam("user_id")

	if err := h.completeHandler.Handle(c.Request().Context(), command.CompleteScheduleCommand{ID: id, UserID: userID}); err != nil {
		return handleAppError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ScheduleHandler) Delete(c echo.Context) error {
	id, _ := uuid.Parse(c.Param("id"))
	userID := c.QueryParam("user_id")

	if err := h.deleteHandler.Handle(c.Request().Context(), command.DeleteScheduleCommand{ID: id, UserID: userID}); err != nil {
		return handleAppError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
