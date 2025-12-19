// Package github provides a client for interacting with the GitHub Actions API.
package github

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/go-github/v67/github"
)

const httpTimeout = 30 * time.Second

// Client wraps the GitHub API client with rate limit logging.
type Client struct {
	client *github.Client
}

// loggingTransport logs rate limit information after each request.
type loggingTransport struct {
	transport http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	rateUsed := resp.Header.Get("X-RateLimit-Used")
	rateLimit := resp.Header.Get("X-RateLimit-Limit")
	if rateUsed != "" && rateLimit != "" {
		slog.Info("request completed",
			"method", req.Method,
			"path", req.URL.Path,
			"status", resp.StatusCode,
			"ratelimit_used", rateUsed,
			"ratelimit_limit", rateLimit)
	}

	return resp, nil
}

// NewClient creates a new GitHub client with the given token.
// If token is empty, requests will be unauthenticated (lower rate limits).
func NewClient(token string) *Client {
	var httpClient *http.Client
	if token != "" {
		httpClient = &http.Client{
			Timeout: httpTimeout,
			Transport: &loggingTransport{
				transport: &github.BasicAuthTransport{
					Username: "x-access-token",
					Password: token,
				},
			},
		}
	} else {
		httpClient = &http.Client{
			Timeout: httpTimeout,
			Transport: &loggingTransport{
				transport: http.DefaultTransport,
			},
		}
	}

	return &Client{
		client: github.NewClient(httpClient),
	}
}

// Actions returns the Actions service for accessing GitHub Actions API.
func (c *Client) Actions() *github.ActionsService {
	return c.client.Actions
}

// RateLimits returns current rate limit status.
func (c *Client) RateLimits(ctx context.Context) (*github.RateLimits, error) {
	limits, _, err := c.client.RateLimit.Get(ctx)
	if err != nil {
		return nil, err
	}
	return limits, nil
}
