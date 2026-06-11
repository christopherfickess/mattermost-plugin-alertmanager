package main

import (
	"strings"
	"testing"
	"time"
)

// TestFormatRotationAge covers the human-label cases for the
// /alertmanager list "Rotated" column. Boundary cases (< 24h →
// "today", < 48h → "yesterday") are exercised explicitly because
// the cutoffs are part of the user-visible contract.
func TestFormatRotationAge(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name    string
		ts      time.Time
		overdue bool
		want    string
	}{
		{"zero value renders never", time.Time{}, false, "never"},
		{"zero value with overdue prefix", time.Time{}, true, "⚠️ never"},
		{"1 hour ago is today", now.Add(-1 * time.Hour), false, "today"},
		{"23 hours ago is today", now.Add(-23 * time.Hour), false, "today"},
		{"25 hours ago is yesterday", now.Add(-25 * time.Hour), false, "yesterday"},
		{"47 hours ago is yesterday", now.Add(-47 * time.Hour), false, "yesterday"},
		{"49 hours ago is 2 days ago", now.Add(-49 * time.Hour), false, "2 days ago"},
		{"100 days ago", now.Add(-100 * 24 * time.Hour), false, "100 days ago"},
		{"overdue prefix on N days ago", now.Add(-100 * 24 * time.Hour), true, "⚠️ 100 days ago"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := formatRotationAge(tc.ts, now, tc.overdue)
			if got != tc.want {
				t.Errorf("formatRotationAge(%v, _, overdue=%v) = %q; want %q",
					tc.ts, tc.overdue, got, tc.want)
			}
		})
	}
}

// TestParseEndToEndExtraLabels covers the admin page's free-text
// label parser. Lenient by design — malformed pairs are dropped
// silently rather than rejecting the form submission.
func TestParseEndToEndExtraLabels(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want map[string]string
	}{
		{"empty input returns nil", "", nil},
		{"single pair", "namespace=billing", map[string]string{"namespace": "billing"}},
		{"multiple pairs", "namespace=billing pod=api-1", map[string]string{"namespace": "billing", "pod": "api-1"}},
		{"no equals dropped", "foo bar", nil},
		{"empty key dropped", "=value", nil},
		{"empty value dropped", "key=", nil},
		{"mixed valid and invalid", "good=ok foo bar=baz", map[string]string{"good": "ok", "bar": "baz"}},
		{"extra whitespace tolerated", "  a=1   b=2  ", map[string]string{"a": "1", "b": "2"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseEndToEndExtraLabels(tc.in)
			if !mapsEqual(got, tc.want) {
				t.Errorf("parseEndToEndExtraLabels(%q) = %v; want %v", tc.in, got, tc.want)
			}
		})
	}
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// Defensive compile check: make sure the rotation age formatter's
// overdue prefix string stays exactly one rune of warning emoji
// plus a space — if someone later changes it, the test catches it
// (rendering width matters for the markdown table column).
func TestRotationAgeOverdueMarker(t *testing.T) {
	got := formatRotationAge(time.Time{}, time.Now(), true)
	if !strings.HasPrefix(got, "⚠️ ") {
		t.Errorf("expected overdue marker prefix, got %q", got)
	}
}
