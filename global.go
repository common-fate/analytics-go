package analytics

import (
	"sync"

	"github.com/common-fate/analytics/acore"
	"go.uber.org/zap"
)

var (
	// globalMu locks concurrent access to the global client.
	globalMu sync.RWMutex

	globalClient = acore.New()
)

func SetDevelopment(dev bool) {
	endpoint := "https://t.commonfate.io"
	if dev {
		endpoint = "https://t-dev.commonfate.io"
	}

	client, err := acore.NewWithConfig(acore.Config{
		Endpoint: endpoint,
	})
	if err != nil {
		zap.L().Named("cf-analytics").Error("error setting client", zap.Error(err))
		return
	}
	ReplaceGlobal(client)
}

// G returns the global client.
func G() *acore.Client {
	globalMu.RLock()
	s := globalClient
	globalMu.RUnlock()
	return s
}

func ReplaceGlobal(c *acore.Client) {
	globalMu.Lock()
	globalClient = c
	globalMu.Unlock()
}

// Close the global client.
func Close() {
	client := G()
	err := client.Close()
	if err != nil {
		zap.L().Named("cf-analytics").Error("error closing client", zap.Error(err))
	}
}
