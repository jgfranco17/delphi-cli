package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextWindow_Bar(t *testing.T) {
	tests := []struct {
		pct      float64
		expected string
	}{
		{0, "░░░░░░░░░░"},
		{100, "▓▓▓▓▓▓▓▓▓▓"},
		{50, "▓▓▓▓▓░░░░░"},
		{10, "▓░░░░░░░░░"},
		{-5, "░░░░░░░░░░"},
		{105, "▓▓▓▓▓▓▓▓▓▓"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("pct=%.0f", tt.pct), func(t *testing.T) {
			cw := &ContextWindow{UsedPercentage: tt.pct}
			assert.Equal(t, tt.expected, cw.Render())
		})
	}
}

func TestContextWindow_Color(t *testing.T) {
	assert.Equal(t, ColorRed, (&ContextWindow{UsedPercentage: 80}).Color())
	assert.Equal(t, ColorRed, (&ContextWindow{UsedPercentage: 100}).Color())
	assert.Equal(t, ColorYellow, (&ContextWindow{UsedPercentage: 50}).Color())
	assert.Equal(t, ColorYellow, (&ContextWindow{UsedPercentage: 79}).Color())
	assert.Equal(t, ColorGreen, (&ContextWindow{UsedPercentage: 0}).Color())
	assert.Equal(t, ColorGreen, (&ContextWindow{UsedPercentage: 49}).Color())
}

func TestCost_Format(t *testing.T) {
	assert.Equal(t, "¥0", (&Cost{}).Format())
	assert.Equal(t, "¥150", (&Cost{TotalCostUSD: 1.0}).Format())
	assert.Equal(t, "¥15", (&Cost{TotalCostUSD: 0.1}).Format())
	assert.Equal(t, "¥8", (&Cost{TotalCostUSD: 0.05}).Format())
}

func TestRateLimits_Format(t *testing.T) {
	assert.Equal(t, "undetermined", (&RateLimits{}).Format())
	assert.Equal(t, "5h: 30%", (&RateLimits{FiveHour: &RateLimit{UsedPercentage: 30}}).Format())
	assert.Equal(t, "7d: 15%", (&RateLimits{SevenDay: &RateLimit{UsedPercentage: 15}}).Format())
	assert.Equal(t, "5h: 30% 7d: 15%", (&RateLimits{
		FiveHour: &RateLimit{UsedPercentage: 30},
		SevenDay: &RateLimit{UsedPercentage: 15},
	}).Format())
}
