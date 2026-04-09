// Package sentinelx provides a Go client for the SentinelX Enforcement API.
// Pre-execution enforcement at the commit boundary.
// https://sentinelx.ai
package sentinelx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultBaseURL = "https://enforce.sentinelx.ai"
	DefaultTimeout = 10 * time.Second
	Version        = "0.1.0"
)

// Client is the SentinelX enforcement client.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// Option configures the client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithTimeout sets the HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

// New creates a new SentinelX client.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		http:    &http.Client{Timeout: DefaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// EnforceRequest is the payload sent to /v2/enforce.
type EnforceRequest struct {
	Action  string         `json:"action"`
	Context map[string]any `json:"context"`
}

// Violation is a single invariant violation.
type Violation struct {
	Primitive  string `json:"primitive"`
	Code       string `json:"code"`
	Constraint string `json:"constraint"`
	Message    string `json:"message"`
}

// Receipt is returned on every enforcement decision.
type Receipt struct {
	Verdict        string      `json:"verdict"`
	Summary        string      `json:"summary"`
	Constraint     *string     `json:"constraint"`
	ConstraintPack string      `json:"constraint_pack"`
	ViolationCode  *string     `json:"violation_code"`
	Violations     []Violation `json:"violations"`
	Mode           string      `json:"mode"`
	EnvelopeClass  string      `json:"envelope_class"`
	TraceID        string      `json:"trace_id"`
	RequestHash    string      `json:"request_hash"`
	ReceiptHash    string      `json:"receipt_hash"`
	InvVersion     string      `json:"inv_version"`
	LatencyMs      int         `json:"latency_ms"`
}

// AdmissibilityError is returned when an action is INADMISSIBLE.
type AdmissibilityError struct {
	Receipt    Receipt
	StatusCode int
}

func (e *AdmissibilityError) Error() string {
	constraint := ""
	if e.Receipt.Constraint != nil {
		constraint = *e.Receipt.Constraint
	}
	return fmt.Sprintf("sentinelx: INADMISSIBLE — %s (%s)", e.Receipt.Summary, constraint)
}

// Enforce evaluates action admissibility at the commit boundary.
// Returns a Receipt on ADMISSIBLE. Returns *AdmissibilityError on INADMISSIBLE.
func (c *Client) Enforce(ctx context.Context, action string, context map[string]any) (*Receipt, error) {
	body, err := json.Marshal(EnforceRequest{Action: action, Context: context})
	if err != nil {
		return nil, fmt.Errorf("sentinelx: marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v2/enforce", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("sentinelx: request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("User-Agent", "sentinelx-go/"+Version)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sentinelx: http error: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sentinelx: read error: %w", err)
	}

	var receipt Receipt
	if err := json.Unmarshal(data, &receipt); err != nil {
		return nil, fmt.Errorf("sentinelx: decode error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &AdmissibilityError{Receipt: receipt, StatusCode: resp.StatusCode}
	}

	return &receipt, nil
}

// Evaluate always returns the receipt without returning an error on INADMISSIBLE.
// Useful for logging pipelines and observe mode.
func (c *Client) Evaluate(ctx context.Context, action string, context map[string]any) (*Receipt, error) {
	receipt, err := c.Enforce(ctx, action, context)
	if err != nil {
		if ae, ok := err.(*AdmissibilityError); ok {
			return &ae.Receipt, nil
		}
		return nil, err
	}
	return receipt, nil
}

// Health checks the API status.
func (c *Client) Health(ctx context.Context) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}
