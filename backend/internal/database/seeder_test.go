package database

import (
	"testing"
	"time"

	"satset-app/internal/model"
)

func TestBuildStatusHistories(t *testing.T) {
	start := time.Date(2026, 10, 18, 2, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		targetStatus string
		wantLen      int
	}{
		{
			name:         "initial status DITERIMA",
			targetStatus: model.WorkStatusDiterima,
			wantLen:      1,
		},
		{
			name:         "mid status DISETRIKA",
			targetStatus: model.WorkStatusDisetrika,
			wantLen:      4,
		},
		{
			name:         "final status SELESAI",
			targetStatus: model.WorkStatusSelesai,
			wantLen:      len(model.ValidWorkStatuses),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			histories := buildStatusHistories(1, tc.targetStatus, start, 2, 3)
			if len(histories) != tc.wantLen {
				t.Fatalf("got %d history entries, want %d", len(histories), tc.wantLen)
			}

			if histories[0].FromStatus != "" || histories[0].ToStatus != model.WorkStatusDiterima {
				t.Fatalf("unexpected first history entry: %+v", histories[0])
			}

			last := histories[len(histories)-1]
			if last.ToStatus != tc.targetStatus {
				t.Fatalf("last status = %q, want %q", last.ToStatus, tc.targetStatus)
			}

			for i := 1; i < len(histories); i++ {
				if !model.IsValidTransition(histories[i].FromStatus, histories[i].ToStatus, false) {
					t.Errorf("invalid transition at index %d: %s -> %s", i, histories[i].FromStatus, histories[i].ToStatus)
				}
				if !histories[i].CreatedAt.After(histories[i-1].CreatedAt) {
					t.Errorf("timestamps not strictly increasing at index %d", i)
				}
			}
		})
	}
}
