package acore

// NoopClient implements acore.Client but does nothing.
type NoopClient struct{}

func (c *NoopClient) Enqueue(Message) error { return nil }

func (c *NoopClient) Close() error { return nil }
