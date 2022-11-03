package analytics

import (
	"context"

	"github.com/common-fate/analytics-go/acore"
)

type contextKey struct {
	name string
}

var analyticsClientContext = contextKey{name: "analyticsClientContext"}

// FromContext loads an analytics client from context.
// analytics.SetContext must have been called.
// if there is no client in context, a new noop client is returned.
func FromContext(ctx context.Context) *Client {
	c, ok := ctx.Value(analyticsClientContext).(*Client)
	if !ok {
		return newClient(&acore.NoopClient{}, defaultLogger())
	}
	return c
}

// SetContext sets the analytics in context
func SetContext(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, analyticsClientContext, client)
}
