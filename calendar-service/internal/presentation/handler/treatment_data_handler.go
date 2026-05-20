package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type TreatmentDataHandler struct {
	adminBaseURL string
	httpClient   *http.Client
}

func NewTreatmentDataHandler(adminBaseURL string) *TreatmentDataHandler {
	return &TreatmentDataHandler{
		adminBaseURL: adminBaseURL,
		httpClient:   &http.Client{Timeout: 3 * time.Second},
	}
}

func (h *TreatmentDataHandler) Register(g *echo.Group) {
	g.GET("/categories", h.proxy)
	g.GET("/categories/:id/treatments", h.proxy)
	g.GET("/treatments/:id/dosage-types", h.proxy)
}

func (h *TreatmentDataHandler) proxy(c echo.Context) error {
	url := h.adminBaseURL + c.Request().URL.Path
	if c.Request().URL.RawQuery != "" {
		url += "?" + c.Request().URL.RawQuery
	}

	req, err := http.NewRequestWithContext(c.Request().Context(), c.Request().Method, url, nil)
	if err != nil {
		return respondError(c, http.StatusBadGateway, "EXTERNAL_API_ERROR", "failed to create proxy request")
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return respondError(c, http.StatusBadGateway, "EXTERNAL_API_ERROR", "admin api unavailable")
	}
	defer resp.Body.Close()

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(resp.StatusCode)
	_, _ = io.Copy(c.Response().Writer, resp.Body)
	return nil
}
