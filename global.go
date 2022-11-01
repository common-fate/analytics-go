package analytics

import (
	"os"
	"strings"
	"sync"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

const (
	DevEndpoint     = "https://t-dev.commonfate.io"
	DefaultEndpoint = "https://t.commonfate.io"
)

var (
	// globalMu locks concurrent access to the global client.
	globalMu sync.RWMutex
	// call analytics.Configure() to set up the client.
	globalClient acore.Client = &acore.NoopClient{}
)

var (
	// Disabled disables analytics altogether.
	Disabled = Config{
		Endpoint: "",
		Enabled:  false,
	}
	// Development uses https://t-dev.commonfate.io as the endpoint.
	Development = Config{
		Endpoint: DevEndpoint,
		Enabled:  true,
	}
	// Default uses https://t.commonfate.io as the endpoint.
	Default = Config{
		Endpoint: DefaultEndpoint,
		Enabled:  true,
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

// Configure the global analytics client.
// Usage:
//
//	analytics.Configure(analytics.Development)
func Configure(c Config) {
	// create a no-op client if analytics are disabled.
	if !c.Enabled {
		c := acore.NoopClient{}
		ReplaceGlobal(&c)
		return
	}

	client, err := acore.NewWithConfig(acore.Config{
		Endpoint: c.Endpoint,
		Callback: debugCallback{},
	})
	if err != nil {
		zap.L().Named("cf-analytics").Error("error setting client", zap.Error(err))
		return
	}
	ReplaceGlobal(client)

	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("configured analytics client", zap.Any("config", c))
	}
}

// ConfigureFromEnv sets up the global analytics client based on the following
// parameters:
//
// - URL is CF_ANALYTICS_URL, or falls back to the default URL if not provided
// - Disabled if CF_ANALYTICS_DISABLED is true
func ConfigureFromEnv() {
	enabled := strings.ToLower(os.Getenv("CF_ANALYTICS_DISABLED")) != "true"

	Configure(Config{
		Endpoint: endpointOrDefault(os.Getenv("CF_ANALYTICS_URL")),
		Enabled:  enabled,
	})
}

// Config is configuration for the global analytics client.
type Config struct {
	Endpoint string
	Enabled  bool
}

// G returns the global client.
func G() acore.Client {
	globalMu.RLock()
	s := globalClient
	globalMu.RUnlock()
	return s
}

func ReplaceGlobal(c acore.Client) {
	globalMu.Lock()
	globalClient = c
	globalMu.Unlock()
}

// Close the global client.
func Close() {
	client := G()
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("closing analytics client", zap.String("url", client.EndpointURL()))
	}

	err := client.Close()
	if err != nil {
		zap.L().Named("cf-analytics").Error("error closing client", zap.Error(err))
	}
}
