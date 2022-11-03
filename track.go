package analytics

import (
	"errors"
	"os"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

type eventMarshaller interface {
	marshalEvent() ([]acore.Message, error)
}

// Track an event using the analytics client.
func (c *Client) Track(e Event) {
	// if the event implements custom marshalling, use it rather than sending a
	// 'Capture' event.
	if em, ok := e.(eventMarshaller); ok {
		events, err := em.marshalEvent()
		if err != nil && os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
			zap.L().Named("cf-analytics").Info("error marshalling analytics events", zap.Error(err))
			return
		}

		for _, evt := range events {
			enqueueAndLog(c.coreclient, evt)
		}

		return
	}

	evt, err := c.marshalToCapture(e)
	if err != nil && os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("error marshalling event", zap.Error(err))
		return
	}

	enqueueAndLog(c.coreclient, evt)
}

// marshalToCapture marshals an Event into an acore.Capture payload to be dispatched.
func (c *Client) marshalToCapture(e Event) (acore.Capture, error) {
	uid := e.userID()
	hashedID, ok := hash(uid, "usr")
	if !ok {
		return acore.Capture{}, errors.New("could not hash user ID")
	}

	typ := e.Type()

	props := eventToProperties(e)

	evt := acore.Capture{
		Event:      typ,
		Properties: props,
		DistinctId: hashedID,
	}

	if c.deploymentID != nil {
		evt.Groups = acore.NewGroups().Set("deployment", *c.deploymentID)
	}

	return evt, nil
}

// enqueueAndLog logs the analytics event using the global zap logger if CF_ANALYTICS_DEBUG is set.
func enqueueAndLog(c acore.Client, m acore.Message) {
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("emitting analytics event", zap.String("url", c.EndpointURL()), zap.Any("event", m))
	}

	c.Enqueue(m)
}
