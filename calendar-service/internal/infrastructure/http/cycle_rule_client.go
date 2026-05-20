package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/model"
	"github.com/NextGenPlatformUnni/calendar-service/internal/domain/port"
	"github.com/google/uuid"
)

type CycleRuleClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCycleRuleClient(baseURL string, timeout time.Duration) *CycleRuleClient {
	return &CycleRuleClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

type cycleRuleResponse struct {
	TreatmentID string `json:"treatment_id"`
	CycleDays   int    `json:"cycle_days"`
}

func (c *CycleRuleClient) GetByTreatmentID(ctx context.Context, treatmentID uuid.UUID) (*port.CycleRule, error) {
	url := fmt.Sprintf("%s/api/cycle-rules/by-treatment/%s", c.baseURL, treatmentID.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, model.NewAppError(model.ErrExternalAPI, "failed to create request", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, model.NewAppError(model.ErrExternalAPI, "admin api call failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, model.NewAppError(model.ErrExternalAPI, fmt.Sprintf("admin api returned %d", resp.StatusCode), nil)
	}

	var result cycleRuleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, model.NewAppError(model.ErrExternalAPI, "failed to decode response", err)
	}

	return &port.CycleRule{TreatmentID: treatmentID, CycleDays: result.CycleDays}, nil
}
