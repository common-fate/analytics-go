package internal

import "testing"

func TestFixturePath(t *testing.T) {
	tests := []struct {
		name  string
		event string
		want  string
	}{
		{
			name:  "ok",
			event: "cf:request.created",
			want:  "fixtures/cf-request-created.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FixturePath(tt.event); got != tt.want {
				t.Errorf("FixturePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
