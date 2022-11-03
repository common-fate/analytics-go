package analytics

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DevEndpoint     = "https://t-dev.commonfate.io"
	DefaultEndpoint = "https://t.commonfate.io"
)

type Client struct {
	mu           *sync.Mutex
	deploymentID *string
	coreclient   acore.Client
	log          *zap.Logger
	// OnFailure is a callback fired if events are failed to be dispatched.
	OnFailure func(e Event)
}

func logLevelFromEnv() zapcore.Level {
	l, err := zapcore.ParseLevel(os.Getenv("CF_ANALYTICS_LOG_LEVEL"))
	if err != nil {
		return zap.PanicLevel
	}
	return l
}

func newClient(coreclient acore.Client, log *zap.Logger) *Client {
	return &Client{
		mu:         &sync.Mutex{},
		coreclient: coreclient,
		log:        log,
	}
}

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

func defaultLogger() *zap.Logger {
	return zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), os.Stderr, logLevelFromEnv())).Named("cf-analytics")
}

// New creates an analytics client.
// Usage:
//
//	analytics.New(analytics.Development)
func New(c Config) *Client {
	log := defaultLogger()

	// create a no-op client if analytics are disabled.
	if !c.Enabled {
		return newClient(&acore.NoopClient{}, log)
	}

	client, err := acore.NewWithConfig(acore.Config{
		Endpoint:  c.Endpoint,
		Verbose:   c.Verbose,
		Interval:  time.Millisecond * 50,
		BatchSize: 3,
	})
	if err != nil {
		log.Error("error setting client", zap.Error(err))
		return newClient(&acore.NoopClient{}, log)
	}

	log.Debug("configured analytics client", zap.Any("config", c))

	return newClient(client, log)
}

// NewFromEnv sets up the analytics client based on the following
// parameters:
//
// - URL is CF_ANALYTICS_URL, or falls back to the default URL if not provided
// - Disabled if CF_ANALYTICS_DISABLED is true
func Env() Config {
	return Config{
		Endpoint: endpointOrDefault(os.Getenv("CF_ANALYTICS_URL")),
		Enabled:  strings.ToLower(os.Getenv("CF_ANALYTICS_DISABLED")) != "true",
		Verbose:  strings.ToLower(os.Getenv("CF_ANALYTICS_LOG_LEVEL")) == "debug",
	}
}

// Config is configuration for the analytics client.
type Config struct {
	Endpoint string `json:"endpoint"`
	Enabled  bool   `json:"enabled"`
	Verbose  bool   `json:"verbose"`
}

// Close the client.
func (c *Client) Close() {
	c.log.Debug("closing analytics client", zap.String("url", c.coreclient.EndpointURL()))

	err := c.coreclient.Close()
	if err != nil {
		c.log.Error("error closing client", zap.Error(err))
	}
}

// SetDeploymentID sets the deployment ID.
func (c *Client) SetDeploymentID(depID string) {
	if depID == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deploymentID = &depID

	c.log.Debug("set deployment", zap.Any("deployment.id", depID))
}
