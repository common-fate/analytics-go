package analytics

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

const (
	DevEndpoint     = "https://t-dev.commonfate.io"
	DefaultEndpoint = "https://t.commonfate.io"
)

type Client struct {
	mu         *sync.Mutex
	deployment *Deployment
	coreclient acore.Client
}

func newClient(coreclient acore.Client) *Client {
	return &Client{
		mu:         &sync.Mutex{},
		coreclient: coreclient,
	}
}

// var (
// 	// globalMu locks concurrent access to the global client.
// 	globalMu sync.RWMutex
// 	// call analytics.Configure() to set up the client.
// 	globalClient Client = &acore.NoopClient{}
// )

var (
	// Disabled disables analytics altogether.
	Disabled = Config{
		Endpoint: "",
		Enabled:  false,
		Verbose:  false,
	}
	// Development uses https://t-dev.commonfate.io as the endpoint.
	Development = Config{
		Endpoint: DevEndpoint,
		Enabled:  true,
		Verbose:  true,
	}
	// Default uses https://t.commonfate.io as the endpoint.
	Default = Config{
		Endpoint: DefaultEndpoint,
		Enabled:  true,
		Verbose:  false,
	}
)

// endpointOrDefault overrides the endpoint or returns the default
// endpoint (https://t.commonfate.io) if the endpoint is empty.
func endpointOrDefault(endpoint string) string {
	if endpoint == "" {
		return DefaultEndpoint
	}
	return endpoint
}

type debugCallback struct{}

func (debugCallback) Success(m acore.Message) {
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("event success", zap.Any("event", m))
	}
}

func (debugCallback) Failure(m acore.Message, err error) {
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Error("event failure", zap.Any("event", m), zap.Error(err))
	}
}

// New the global analytics client.
// Usage:
//
//	analytics.New(analytics.Development)
func New(c Config) *Client {
	// create a no-op client if analytics are disabled.
	if !c.Enabled {
		return newClient(&acore.NoopClient{})
	}

	client, err := acore.NewWithConfig(acore.Config{
		Endpoint:  c.Endpoint,
		Callback:  debugCallback{},
		Verbose:   c.Verbose,
		Interval:  time.Millisecond * 50,
		BatchSize: 2,
	})
	if err != nil {
		zap.L().Named("cf-analytics").Error("error setting client", zap.Error(err))
		return newClient(&acore.NoopClient{})
	}

	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("configured analytics client", zap.Any("config", c))
	}

	return newClient(client)
}

// NewFromEnv sets up the global analytics client based on the following
// parameters:
//
// - URL is CF_ANALYTICS_URL, or falls back to the default URL if not provided
// - Disabled if CF_ANALYTICS_DISABLED is true
func NewFromEnv() *Client {
	return New(Config{
		Endpoint: endpointOrDefault(os.Getenv("CF_ANALYTICS_URL")),
		Enabled:  strings.ToLower(os.Getenv("CF_ANALYTICS_DISABLED")) != "true",
		Verbose:  strings.ToLower(os.Getenv("CF_ANALYTICS_DEBUG")) == "true",
	})
}

// Config is configuration for the global analytics client.
type Config struct {
	Endpoint string `json:"endpoint"`
	Enabled  bool   `json:"enabled"`
	Verbose  bool   `json:"verbose"`
}

// // G returns the global client.
// func G() acore.Client {
// 	globalMu.RLock()
// 	s := globalClient
// 	globalMu.RUnlock()
// 	return s
// }

// func ReplaceGlobal(c acore.Client) {
// 	globalMu.Lock()
// 	globalClient = c
// 	globalMu.Unlock()
// }

// Close the global client.
func (c *Client) Close() {
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("closing analytics client", zap.String("url", c.coreclient.EndpointURL()))
	}

	err := c.coreclient.Close()
	if err != nil {
		zap.L().Named("cf-analytics").Error("error closing client", zap.Error(err))
	}
}
