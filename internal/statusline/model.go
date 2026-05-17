package statusline

import (
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
)

const conversionRateUSDJPY = 150.0

// Input is the JSON structure piped from Claude Code.
type Input struct {
	Model         Model         `json:"model"`
	Workspace     Workspace     `json:"workspace"`
	ContextWindow ContextWindow `json:"context_window"`
	Cost          Cost          `json:"cost"`
	RateLimits    RateLimits    `json:"rate_limits"`
}

type Model struct {
	DisplayName string `json:"display_name"`
}

type Workspace struct {
	CurrentDir string `json:"current_dir"`
}

type ContextWindow struct {
	UsedPercentage float64 `json:"used_percentage"`
}

// Render renders a fixed-width block character bar for the window's usage percentage.
func (c *ContextWindow) Render() string {
	pct := int(c.UsedPercentage)
	switch {
	case pct < 0:
		pct = 0
	case pct > 100:
		pct = 100
	}
	filled := pct * barWidth / 100
	empty := barWidth - filled
	return strings.Repeat("▓", filled) + strings.Repeat("░", empty)
}

// Color returns the color appropriate for the window's usage percentage.
func (c *ContextWindow) Color() *color.Color {
	pct := int(c.UsedPercentage)
	switch {
	case pct >= 80:
		return colorRed
	case pct >= 50:
		return colorYellow
	default:
		return colorGreen
	}
}

type Cost struct {
	TotalCostUSD float64 `json:"total_cost_usd"`
}

// Format converts the USD cost to a JPY string at a fixed 150 rate.
func (c *Cost) Format() string {
	jpy := math.Round(c.TotalCostUSD * conversionRateUSDJPY)
	return fmt.Sprintf("¥%.0f", jpy)
}

type RateLimits struct {
	FiveHour *RateLimit `json:"five_hour"`
	SevenDay *RateLimit `json:"seven_day"`
}

// Format returns a human-readable string of the rate limit percentages.
// Returns "undetermined" when both limits are absent.
func (r *RateLimits) Format() string {
	if r.FiveHour == nil && r.SevenDay == nil {
		return "undetermined"
	}
	parts := make([]string, 0, 2)
	if r.FiveHour != nil {
		parts = append(parts, fmt.Sprintf("5h: %.0f%%", r.FiveHour.UsedPercentage))
	}
	if r.SevenDay != nil {
		parts = append(parts, fmt.Sprintf("7d: %.0f%%", r.SevenDay.UsedPercentage))
	}
	return strings.Join(parts, " ")
}

type RateLimit struct {
	UsedPercentage float64 `json:"used_percentage"`
}
