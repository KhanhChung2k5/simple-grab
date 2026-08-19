package service

import (
	"testing"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
)

func TestIsValidTransition(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want bool
	}{
		{
			name: "accepted to in progress",
			from: string(model.RideAccepted),
			to:   string(model.RideInProgress),
			want: true,
		},
		{
			name: "in progress to completed",
			from: string(model.RideInProgress),
			to:   string(model.RideCompleted),
			want: true,
		},
		{
			name: "pending to completed",
			from: string(model.RidePending),
			to:   string(model.RideCompleted),
			want: false,
		},
		{
			name: "accepted to completed",
			from: string(model.RideAccepted),
			to:   string(model.RideCompleted),
			want: false,
		},
		{
			name: "completed to in progress",
			from: string(model.RideCompleted),
			to:   string(model.RideInProgress),
			want: false,
		},
		{
			name: "accepted to cancelled",
			from: string(model.RideAccepted),
			to:   string(model.RideCancelled),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidTransition(tt.from, tt.to)

			if got != tt.want {
				t.Fatalf(
					"isValidTransition(%q, %q) = %v; want %v",
					tt.from,
					tt.to,
					got,
					tt.want,
				)
			}
		})
	}
}
