package service

import "testing"

func TestNormalizeAdminStatsDateKeyHandlesRFC3339(t *testing.T) {
	got := normalizeAdminStatsDateKey("2026-07-25T00:00:00Z")
	if got != "2026-07-25" {
		t.Fatalf("expected RFC3339 date to normalize to 2026-07-25, got %q", got)
	}
}

func TestNormalizeAdminStatsDateKeyKeepsDateOnlyValue(t *testing.T) {
	got := normalizeAdminStatsDateKey("2026-07-25")
	if got != "2026-07-25" {
		t.Fatalf("expected date-only value to stay unchanged, got %q", got)
	}
}
