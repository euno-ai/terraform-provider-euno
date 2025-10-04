package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// EunoClient represents the API client for Euno
type EunoClient struct {
	serverURL    string
	apiKey       string
	accountID    int
	httpClient   *http.Client
	rateLimiter  chan struct{}
}

// NewEunoClient creates a new Euno API client with rate limiting
func NewEunoClient(serverURL, apiKey string, accountID int) *EunoClient {
	return &EunoClient{
		serverURL:   serverURL,
		apiKey:      apiKey,
		accountID:   accountID,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		rateLimiter: make(chan struct{}, 3), // Allow max 3 concurrent requests
	}
}

// acquireRateLimit acquires a slot in the rate limiter
func (c *EunoClient) acquireRateLimit(ctx context.Context) error {
	select {
	case c.rateLimiter <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// releaseRateLimit releases a slot in the rate limiter
func (c *EunoClient) releaseRateLimit() {
	<-c.rateLimiter
}

// IntegrationSchedule represents the schedule configuration
type IntegrationSchedule struct {
	TimeZone    string   `json:"time_zone"`
	RepeatOn    []string `json:"repeat_on,omitempty"`
	RepeatTime  string   `json:"repeat_time,omitempty"`
	RepeatPeriod *int    `json:"repeat_period,omitempty"`
}

// InvalidationStrategy represents the invalidation strategy configuration
type InvalidationStrategy struct {
	RevisionID *int `json:"revision_id,omitempty"`
	TTLDays    int  `json:"ttl_days"`
}

// IntegrationIn represents the input for creating an integration
type IntegrationIn struct {
	IntegrationType string               `json:"integration_type"`
	Name            string               `json:"name"`
	Active          bool                 `json:"active"`
	Schedule        *IntegrationSchedule `json:"schedule,omitempty"`
	Configuration   map[string]interface{} `json:"configuration"`
	InvalidationStrategy *InvalidationStrategy `json:"invalidation_strategy,omitempty"`
	PendingCredentialsLookupKey *string  `json:"pending_credentials_lookup_key,omitempty"`
}

// IntegrationOut represents the output from the API
type IntegrationOut struct {
	ID                        int                    `json:"id"`
	IntegrationType           string                 `json:"integration_type"`
	AccountID                 int                    `json:"account_id"`
	CreatedAt                 string                 `json:"created_at"`
	CreatedBy                 string                 `json:"created_by"`
	LastUpdatedAt             string                 `json:"last_updated_at"`
	LastUpdatedBy             string                 `json:"last_updated_by"`
	Name                      string                 `json:"name"`
	Active                    *bool                  `json:"active"`
	Configuration             map[string]interface{} `json:"configuration"`
	Schedule                  *IntegrationSchedule   `json:"schedule"`
	CollectedIntegrationData  map[string]interface{} `json:"collected_integration_data"`
	LastRunStatus             *string                 `json:"last_run_status"`
	LastCompletedRunEndTime    *string                `json:"last_completed_run_end_time"`
	Health                    *string                 `json:"health"`
	TriggerType               *string                `json:"trigger_type"`
	TriggerSecret             *string                `json:"trigger_secret"`
	TriggerURL                *string                `json:"trigger_url"`
	InvalidationStrategy      *InvalidationStrategy  `json:"invalidation_strategy"`
	LastTimeTriggered         *string                `json:"last_time_triggered"`
	PendingCredentialsLookupKey *string              `json:"pending_credentials_lookup_key"`
}

// CreateIntegration creates a new integration
func (c *EunoClient) CreateIntegration(ctx context.Context, integration IntegrationIn) (*IntegrationOut, error) {
	if err := c.acquireRateLimit(ctx); err != nil {
		return nil, err
	}
	defer c.releaseRateLimit()

	url := fmt.Sprintf("%s/accounts/%d/integrations", c.serverURL, c.accountID)
	
	jsonData, err := json.Marshal(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal integration data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result IntegrationOut
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetIntegration retrieves an integration by ID
func (c *EunoClient) GetIntegration(ctx context.Context, integrationID int) (*IntegrationOut, error) {
	if err := c.acquireRateLimit(ctx); err != nil {
		return nil, err
	}
	defer c.releaseRateLimit()

	url := fmt.Sprintf("%s/accounts/%d/integrations/%d", c.serverURL, c.accountID, integrationID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("integration not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result IntegrationOut
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateIntegration updates an existing integration
func (c *EunoClient) UpdateIntegration(ctx context.Context, integrationID int, integration IntegrationIn) (*IntegrationOut, error) {
	if err := c.acquireRateLimit(ctx); err != nil {
		return nil, err
	}
	defer c.releaseRateLimit()

	url := fmt.Sprintf("%s/accounts/%d/integrations/%d", c.serverURL, c.accountID, integrationID)
	
	jsonData, err := json.Marshal(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal integration data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("integration not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result IntegrationOut
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteIntegration deletes an integration
func (c *EunoClient) DeleteIntegration(ctx context.Context, integrationID int) error {
	if err := c.acquireRateLimit(ctx); err != nil {
		return err
	}
	defer c.releaseRateLimit()

	url := fmt.Sprintf("%s/accounts/%d/integrations/%d", c.serverURL, c.accountID, integrationID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("integration not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
