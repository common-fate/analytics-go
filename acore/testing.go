package acore

import (
	"testing"
	"time"
)

type TestConfig struct {
	Now func() time.Time
	UID func() string
}

func NewTestWithConfig(t *testing.T, config Config, tc TestConfig) Client {
	config.now = tc.Now
	config.uid = tc.UID
	c, err := NewWithConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
