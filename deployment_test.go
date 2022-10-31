package analytics

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_globalDep_loadDeployment(t *testing.T) {
	type fields struct {
		mu      *sync.Mutex
		dep     *Deployment
		timeout time.Duration
		loader  func(ctx context.Context) (*Deployment, error)
	}
	tests := []struct {
		name   string
		fields fields
		want   *Deployment
	}{
		{
			name: "ok",
			fields: fields{
				mu: &sync.Mutex{},
				loader: func(ctx context.Context) (*Deployment, error) {
					return &Deployment{ID: "test"}, nil
				},
			},
			want: &Deployment{
				ID: "dep_TZZ6MBEb8p8OugHESLN1wWKbL-0BzfzDrtkfG1fV3V4",
			},
		},
		{
			name: "load with no id should not be hashed",
			fields: fields{
				mu: &sync.Mutex{},
				loader: func(ctx context.Context) (*Deployment, error) {
					return &Deployment{ID: ""}, nil
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := &globalDep{
				mu:      tt.fields.mu,
				dep:     tt.fields.dep,
				timeout: tt.fields.timeout,
				loader:  tt.fields.loader,
			}
			ctx := context.Background()
			gd.loadDeployment(ctx)
			got := gd.Get(ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}
