package analytics

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testLoader struct {
	Deployment *Deployment
}

func (tl testLoader) LoadDeployment(ctx context.Context) (*Deployment, error) {
	return tl.Deployment, nil
}

func Test_globalDep_loadDeployment(t *testing.T) {
	type fields struct {
		mu      *sync.Mutex
		dep     *Deployment
		timeout time.Duration
		loader  DeploymentLoader
	}
	tests := []struct {
		name   string
		fields fields
		want   *Deployment
	}{
		{
			name: "ok",
			fields: fields{
				mu:     &sync.Mutex{},
				loader: &testLoader{Deployment: &Deployment{ID: "test"}},
			},
			want: &Deployment{
				ID: "test",
			},
		},
		{
			name: "load with no id return nil",
			fields: fields{
				mu:     &sync.Mutex{},
				loader: &testLoader{Deployment: &Deployment{ID: ""}},
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
