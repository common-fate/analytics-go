package analytics

import (
	"os"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

// Track an event using the global analytics client.
func (c *Client) Track(e Event) {
	e = hashValues(e)
	uid := e.userID()
	typ := e.Type()

	evt := acore.Track{
		Event:      typ,
		Properties: e,
		UserId:     uid,
	}

	// generate an anonymous ID if there is no user ID.
	if uid == "" {
		evt.AnonymousId = c.uid()
	}

	if c.deployment != nil {
		enqueueAndLog(c.coreclient, acore.Group{
			GroupId:     c.deployment.ID,
			Traits:      c.deployment.Traits(),
			UserId:      evt.UserId,
			AnonymousId: evt.AnonymousId, // one of UserId or AnonymousId will be set.
		})
	}

	enqueueAndLog(c.coreclient, evt)
}

// enqueueAndLog logs the analytics event using the global zap logger if CF_ANALYTICS_DEBUG is set.
func enqueueAndLog(c acore.Client, m acore.Message) {
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("emitting analytics event", zap.String("url", c.EndpointURL()), zap.Any("event", m))
	}

	c.Enqueue(m)
}
