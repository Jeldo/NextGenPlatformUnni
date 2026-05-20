package handler

import (
	"net/http"

	"github.com/NextGenPlatformUnni/calendar-service/internal/application/command"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/labstack/echo/v4"
)

type MockHandler struct {
	reservationHandler *command.HandleReservationFixedCommandHandler
}

func NewMockHandler(reservationHandler *command.HandleReservationFixedCommandHandler) *MockHandler {
	return &MockHandler{reservationHandler: reservationHandler}
}

func (h *MockHandler) Register(g *echo.Group) {
	g.POST("/events/reservation-fixed", h.HandleReservationFixed)
}

func (h *MockHandler) HandleReservationFixed(c echo.Context) error {
	var event model.ReservationFixedEvent
	if err := c.Bind(&event); err != nil {
		return respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid event payload")
	}

	if err := h.reservationHandler.Handle(c.Request().Context(), command.HandleReservationFixedCommand{Event: event}); err != nil {
		return handleAppError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "processed"})
}
