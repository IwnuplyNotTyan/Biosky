package bluesky

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	profileCollection = "app.bsky.actor.profile"
	profileRKey       = "self"
	maxResponseBytes  = 1 << 20
)

// Config holds the credentials needed to talk to a PDS.
type Config struct {
	Handle      string
	AppPassword string
	PDS         string
	UserAgent   string
}

// Client is a minimal AT Protocol client for one profile.
type Client struct {
	cfg       Config
	pds       string
	http      *http.Client
	accessJWT string
	did       string
}

// APIError describes a non-2xx XRPC response.
type APIError struct {
	Status  int
	Code    string
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Status, e.Code, e.Message)
}

// New returns a client with a normalized PDS and a 15-second timeout.
func New(cfg Config) *Client {
	pds := strings.TrimRight(strings.TrimSpace(cfg.PDS), "/")
	if pds == "" {
		pds = "https://bsky.social"
	}
	if !strings.Contains(pds, "://") {
		pds = "https://" + pds
	}
	ua := cfg.UserAgent
	if ua == "" {
		ua = "biosky/dev"
	}
	cfg.PDS = pds
	cfg.UserAgent = ua
	return &Client{
		cfg:  cfg,
		pds:  pds,
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

// DID returns the authenticated account DID.
func (c *Client) DID() string { return c.did }

// Login creates a session from the app password.
func (c *Client) Login(ctx context.Context) error {
	in := map[string]string{
		"identifier": c.cfg.Handle,
		"password":   c.cfg.AppPassword,
	}
	var out struct {
		AccessJWT string `json:"accessJwt"`
		DID       string `json:"did"`
	}
	if err := c.do(ctx, http.MethodPost, "com.atproto.server.createSession", nil, in, &out); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	if out.AccessJWT == "" || out.DID == "" {
		return errors.New("create session: empty token or did")
	}
	c.accessJWT = out.AccessJWT
	c.did = out.DID
	return nil
}

// UpdateDescription replaces only the description and reports whether it wrote.
func (c *Client) UpdateDescription(ctx context.Context, description string) (bool, error) {
	if c.accessJWT == "" || c.did == "" {
		return false, errors.New("not logged in")
	}
	p, err := c.profile(ctx)
	if err != nil {
		return false, err
	}
	if p.Exists && p.Description() == description {
		return false, nil
	}
	if err := p.SetDescription(description); err != nil {
		return false, err
	}
	if err := c.putProfile(ctx, p); err != nil {
		return false, err
	}
	return true, nil
}

type profile struct {
	Fields map[string]json.RawMessage
	CID    string
	Exists bool
}

// Description decodes the current description, or returns an empty string.
func (p profile) Description() string {
	raw, ok := p.Fields["description"]
	if !ok {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) != nil {
		return ""
	}
	return s
}

// SetDescription replaces only the description and keeps every other field.
func (p *profile) SetDescription(s string) error {
	raw, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal description: %w", err)
	}
	if p.Fields == nil {
		p.Fields = map[string]json.RawMessage{}
	}
	p.Fields["description"] = raw
	if _, ok := p.Fields["$type"]; !ok {
		p.Fields["$type"] = json.RawMessage(`"app.bsky.actor.profile"`)
	}
	return nil
}

func (c *Client) profile(ctx context.Context) (profile, error) {
	query := url.Values{}
	query.Set("repo", c.did)
	query.Set("collection", profileCollection)
	query.Set("rkey", profileRKey)
	var out struct {
		CID   string                     `json:"cid"`
		Value map[string]json.RawMessage `json:"value"`
	}
	if err := c.do(ctx, http.MethodGet, "com.atproto.repo.getRecord", query, nil, &out); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Code == "RecordNotFound" {
			return profile{Fields: map[string]json.RawMessage{}, Exists: false}, nil
		}
		return profile{}, err
	}
	return profile{Fields: out.Value, CID: out.CID, Exists: true}, nil
}

func (c *Client) putProfile(ctx context.Context, p profile) error {
	in := struct {
		Repo       string                     `json:"repo"`
		Collection string                     `json:"collection"`
		RKey       string                     `json:"rkey"`
		Record     map[string]json.RawMessage `json:"record"`
		SwapRecord string                     `json:"swapRecord,omitempty"`
	}{
		Repo:       c.did,
		Collection: profileCollection,
		RKey:       profileRKey,
		Record:     p.Fields,
		SwapRecord: p.CID,
	}
	return c.do(ctx, http.MethodPost, "com.atproto.repo.putRecord", nil, in, nil)
}

func (c *Client) do(ctx context.Context, method, xrpc string, query url.Values, in, out any) error {
	endpoint := c.pds + "/xrpc/" + xrpc
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}
	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.accessJWT != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessJWT)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, xrpc, err)
	}
	defer func() { _ = resp.Body.Close() }()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		apiErr := &APIError{Status: resp.StatusCode}
		var payload struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if json.Unmarshal(data, &payload) == nil && payload.Error != "" {
			apiErr.Code = payload.Error
			apiErr.Message = payload.Message
		} else {
			apiErr.Message = truncate(string(data), 512)
		}
		return apiErr
	}
	if out != nil {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
