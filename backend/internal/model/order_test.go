package model_test

import (
	"testing"

	"satset-app/internal/model"
)

func TestWorkStatusTransitions(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		isOwner bool
		want    bool
	}{
		// Standard flow
		{"Diterima to Dicuci", model.WorkStatusDiterima, model.WorkStatusDicuci, false, true},
		{"Dicuci to Dikeringkan", model.WorkStatusDicuci, model.WorkStatusDikeringkan, false, true},
		{"Dikeringkan to Disetrika", model.WorkStatusDikeringkan, model.WorkStatusDisetrika, false, true},
		{"Disetrika to Dikemas", model.WorkStatusDisetrika, model.WorkStatusDikemas, false, true},
		{"Dikemas to SiapDiambil", model.WorkStatusDikemas, model.WorkStatusSiapDiambil, false, true},
		{"SiapDiambil to Selesai", model.WorkStatusSiapDiambil, model.WorkStatusSelesai, false, true},

		// Valid forward bypasses (e.g. cuci kering tanpa setrika)
		{"Dicuci to Dikemas", model.WorkStatusDicuci, model.WorkStatusDikemas, false, true},
		{"Diterima to Disetrika (setrika saja)", model.WorkStatusDiterima, model.WorkStatusDisetrika, false, true},

		// Illegal transitions for non-owner
		{"Diterima to Selesai (illegal skip)", model.WorkStatusDiterima, model.WorkStatusSelesai, false, false},
		{"Selesai to Dicuci (terminal)", model.WorkStatusSelesai, model.WorkStatusDicuci, false, false},
		{"Dibatalkan to Diterima (terminal)", model.WorkStatusDibatalkan, model.WorkStatusDiterima, false, false},
		{"Selesai to Dibatalkan (cannot cancel completed)", model.WorkStatusSelesai, model.WorkStatusDibatalkan, true, false},

		// Cancellation rules
		{"Diterima to Dibatalkan by Kasir", model.WorkStatusDiterima, model.WorkStatusDibatalkan, false, true},
		{"Dicuci to Dibatalkan by non-owner (denied)", model.WorkStatusDicuci, model.WorkStatusDibatalkan, false, false},
		{"Dicuci to Dibatalkan by Owner (allowed)", model.WorkStatusDicuci, model.WorkStatusDibatalkan, true, true},

		// Owner override
		{"Owner override backward transition", model.WorkStatusDikemas, model.WorkStatusDicuci, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.IsValidTransition(tt.from, tt.to, tt.isOwner)
			if got != tt.want {
				t.Errorf("IsValidTransition(%s, %s, %v) = %v; want %v", tt.from, tt.to, tt.isOwner, got, tt.want)
			}
		})
	}
}

func TestTerminalStatus(t *testing.T) {
	if !model.IsTerminalStatus(model.WorkStatusSelesai) {
		t.Errorf("SELESAI should be terminal")
	}
	if !model.IsTerminalStatus(model.WorkStatusDibatalkan) {
		t.Errorf("DIBATALKAN should be terminal")
	}
	if model.IsTerminalStatus(model.WorkStatusDiterima) {
		t.Errorf("DITERIMA should not be terminal")
	}
	if model.IsTerminalStatus(model.WorkStatusDicuci) {
		t.Errorf("DICUCI should not be terminal")
	}
}
