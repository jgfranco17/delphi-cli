package render

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jgfranco17/delphi-cli/internal/tooling/git"
	"github.com/jgfranco17/delphi-cli/internal/tooling/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGitProvider struct {
	branch string
	dirty  bool
	err    error
}

func (m *mockGitProvider) BranchStatus(_ context.Context, _ string) (string, bool, error) {
	return m.branch, m.dirty, m.err
}

var _ git.Provider = (*mockGitProvider)(nil)

func newTestInput() *model.AgentInput {
	return &model.AgentInput{
		Model:         model.Model{DisplayName: "claude-opus-4-5"},
		Workspace:     model.Workspace{CurrentDir: "/home/user/project"},
		ContextWindow: model.ContextWindow{UsedPercentage: 42},
		Cost:          model.Cost{TotalCostUSD: 1.0},
		RateLimits: model.RateLimits{
			FiveHour: &model.RateLimit{UsedPercentage: 30},
			SevenDay: &model.RateLimit{UsedPercentage: 15},
		},
	}
}

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name        string
		git         *mockGitProvider
		opts        Options
		mutate      func(*model.AgentInput)
		contains    []string
		notContains []string
	}{
		{
			name: "all options enabled",
			git:  &mockGitProvider{branch: "main", dirty: false},
			opts: Options{WithCurrentDir: true, WithGitInfo: true, WithUsageStats: true},
			contains: []string{
				"claude-opus-4-5", "/home/user/project",
				"42%", "¥150", "5h: 30%", "7d: 15%",
				"main", "✓ clean",
			},
		},
		{
			name:        "no options — model only",
			git:         &mockGitProvider{branch: "main"},
			opts:        Options{},
			contains:    []string{"claude-opus-4-5"},
			notContains: []string{"/home/user/project", "42%", "main"},
		},
		{
			name:        "only current dir",
			git:         &mockGitProvider{branch: "main"},
			opts:        Options{WithCurrentDir: true},
			contains:    []string{"claude-opus-4-5", "/home/user/project"},
			notContains: []string{"42%", "main"},
		},
		{
			name:        "only usage stats",
			git:         &mockGitProvider{branch: "main"},
			opts:        Options{WithUsageStats: true},
			contains:    []string{"42%", "¥150", "5h: 30%", "7d: 15%"},
			notContains: []string{"/home/user/project", "main"},
		},
		{
			name:        "only git info",
			git:         &mockGitProvider{branch: "main"},
			opts:        Options{WithGitInfo: true},
			contains:    []string{"main", "✓ clean"},
			notContains: []string{"/home/user/project", "42%"},
		},
		{
			name:     "dirty branch",
			git:      &mockGitProvider{branch: "feature/foo", dirty: true},
			opts:     Options{WithGitInfo: true},
			contains: []string{"feature/foo", "✗ dirty"},
		},
		{
			name:     "no git repo",
			git:      &mockGitProvider{err: errors.New("not a repo")},
			opts:     Options{WithGitInfo: true},
			contains: []string{"none"},
		},
		{
			name: "no rate limits",
			git:  &mockGitProvider{branch: "main"},
			opts: Options{WithUsageStats: true},
			mutate: func(in *model.AgentInput) {
				in.RateLimits.FiveHour = nil
				in.RateLimits.SevenDay = nil
			},
			contains: []string{"undetermined"},
		},
		{
			name: "high context usage",
			git:  &mockGitProvider{branch: "main"},
			opts: Options{WithUsageStats: true},
			mutate: func(in *model.AgentInput) {
				in.ContextWindow.UsedPercentage = 85
			},
			contains: []string{"85%"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := newTestInput()
			if tc.mutate != nil {
				tc.mutate(input)
			}
			var buf bytes.Buffer
			r := newRenderer(&buf, tc.git, tc.opts)

			require.NoError(t, r.Render(t.Context(), input))

			out := buf.String()
			for _, want := range tc.contains {
				assert.Contains(t, out, want)
			}
			for _, absent := range tc.notContains {
				assert.NotContains(t, out, absent)
			}
		})
	}
}

func TestRenderer_GenerateFrom(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		opts     Options
		wantErr  bool
		contains []string
	}{
		{
			name: "valid JSON",
			json: `{
				"model": {"display_name": "claude-mock"},
				"workspace": {"current_dir": "/tmp"},
				"context_window": {"used_percentage": 10},
				"cost": {"total_cost_usd": 0.5},
				"rate_limits": {}
			}`,
			opts:     Options{WithCurrentDir: true},
			contains: []string{"claude-mock", "/tmp"},
		},
		{
			name:    "invalid JSON",
			json:    "{not json}",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			r := newRenderer(&buf, &mockGitProvider{branch: "main"}, tc.opts)

			err := r.GenerateFrom(t.Context(), strings.NewReader(tc.json))
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			out := buf.String()
			for _, want := range tc.contains {
				assert.Contains(t, out, want)
			}
		})
	}
}
