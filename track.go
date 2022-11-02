package analytics

import (
	"os"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

type payloader interface {
	payloads() []acore.Message
}

type deploymentEventer interface {
	deploymentEvent() bool
}

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

	// if true, don't identify the user.
	var isDeploymentEvent bool

	if d, ok := e.(deploymentEventer); ok {
		isDeploymentEvent = d.deploymentEvent()
	}

	if uid != "" && !isDeploymentEvent {
		enqueueAndLog(c.coreclient, acore.Identify{
			UserId: uid,
			Traits: acore.NewTraits().Set("user_id", uid),
		})
	}

	if pl, ok := e.(payloader); ok {
		for _, m := range pl.payloads() {
			enqueueAndLog(c.coreclient, m)
		}
	}
}

// enqueueAndLog logs the analytics event using the global zap logger if CF_ANALYTICS_DEBUG is set.
func enqueueAndLog(c acore.Client, m acore.Message) {
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("emitting analytics event", zap.String("url", c.EndpointURL()), zap.Any("event", m))
	}

	c.Enqueue(m)
}
