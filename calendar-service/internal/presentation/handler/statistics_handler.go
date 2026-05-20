package handler

import (
	"net/http"

	"github.com/NextGenPlatformUnni/calendar-service/internal/application/query"
	"github.com/NextGenPlatformUnni/calendar-service/internal/presentation/dto"
	"github.com/labstack/echo/v4"
)

type StatisticsHandler struct {
	handler *query.GetStatisticsQueryHandler
}

func NewStatisticsHandler(handler *query.GetStatisticsQueryHandler) *StatisticsHandler {
	return &StatisticsHandler{handler: handler}
}

func (h *StatisticsHandler) Register(g *echo.Group) {
	g.GET("", h.Get)
}

func (h *StatisticsHandler) Get(c echo.Context) error {
	userID := c.QueryParam("user_id")

	stats, err := h.handler.Handle(c.Request().Context(), query.GetStatisticsQuery{UserID: userID})
	if err != nil {
		return handleAppError(c, err)
	}

	items := make([]dto.StatItem, len(stats))
	for i, s := range stats {
		items[i] = dto.StatItem{
			TreatmentID:   s.TreatmentID.String(),
			TreatmentName: s.TreatmentName,
			CategoryName:  s.CategoryName,
			Count:         s.Count,
		}
	}

	return c.JSON(http.StatusOK, dto.StatisticsResponse{Stats: items})
}
