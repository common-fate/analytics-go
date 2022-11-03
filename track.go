package analytics

import (
	"errors"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

type eventMarshaller interface {
	marshalEvent(ctx marshalContext) ([]acore.Message, error)
}

// Track an event using the analytics client.
func (c *Client) Track(e Event) {
	// if the event implements custom marshalling, use it rather than sending a
	// 'Capture' event.
	if em, ok := e.(eventMarshaller); ok {
		events, err := em.marshalEvent(marshalContext{DeploymentID: c.deploymentID})
		if err != nil {
			c.log.Error("error marshalling analytics events", zap.Error(err))
			c.failed(e)
			return
		}

		for _, evt := range events {
			err := c.enqueueAndLog(evt)
			if err != nil {
				c.failed(e)
			}
		}

		return
	}

	evt, err := c.marshalToCapture(e)
	if err != nil {
		c.log.Error("error marshalling event", zap.Error(err))
		c.failed(e)
		return
	}

	err = c.enqueueAndLog(evt)
	if err != nil {
		c.failed(e)
	}
}

func (c *Client) failed(e Event) {
	if c.OnFailure != nil {
		c.OnFailure(e)
	}
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

// enqueueAndLog logs the analytics event using the logger if CF_ANALYTICS_LOG_LEVEL is set to 'debug'.
func (c *Client) enqueueAndLog(m acore.Message) error {
	c.log.Debug("emitting analytics event", zap.String("url", c.coreclient.EndpointURL()), zap.Any("event", m))
	return c.coreclient.Enqueue(m)
}
