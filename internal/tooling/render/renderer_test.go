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

func TestRenderer_Render_ContainsExpectedText(t *testing.T) {
	git := &mockGitProvider{branch: "main", dirty: false}
	var buf bytes.Buffer
	r := NewRenderer(&buf, git)

	require.NoError(t, r.Render(t.Context(), newTestInput()))

	out := buf.String()
	assert.Contains(t, out, "claude-opus-4-5")
	assert.Contains(t, out, "/home/user/project")
	assert.Contains(t, out, "42%")
	assert.Contains(t, out, "¥150")
	assert.Contains(t, out, "5h: 30%")
	assert.Contains(t, out, "7d: 15%")
	assert.Contains(t, out, "main")
	assert.Contains(t, out, "✓ clean")
}

func TestRenderer_Render_DirtyBranch(t *testing.T) {
	git := &mockGitProvider{branch: "feature/foo", dirty: true}
	var buf bytes.Buffer
	r := NewRenderer(&buf, git)

	require.NoError(t, r.Render(t.Context(), newTestInput()))

	out := buf.String()
	assert.Contains(t, out, "feature/foo")
	assert.Contains(t, out, "✗ dirty")
}

func TestRenderer_Render_NoGitRepo(t *testing.T) {
	git := &mockGitProvider{err: errors.New("not a repo")}
	var buf bytes.Buffer
	r := NewRenderer(&buf, git)

	require.NoError(t, r.Render(t.Context(), newTestInput()))
	assert.Contains(t, buf.String(), "none")
}

func TestRenderer_Render_NoRateLimits(t *testing.T) {
	git := &mockGitProvider{branch: "main"}
	var buf bytes.Buffer
	r := NewRenderer(&buf, git)

	input := newTestInput()
	input.RateLimits.FiveHour = nil
	input.RateLimits.SevenDay = nil

	require.NoError(t, r.Render(t.Context(), input))
	assert.Contains(t, buf.String(), "undetermined")
}

func TestRenderer_Render_HighUsageBar(t *testing.T) {
	git := &mockGitProvider{branch: "main"}
	var buf bytes.Buffer
	r := NewRenderer(&buf, git)

	input := newTestInput()
	input.ContextWindow.UsedPercentage = 85

	require.NoError(t, r.Render(t.Context(), input))
	assert.Contains(t, buf.String(), "85%")
}

func TestRenderer_RenderFromReader_ValidJSON(t *testing.T) {
	json := `{
		"model": {"display_name": "claude-3"},
		"workspace": {"current_dir": "/tmp"},
		"context_window": {"used_percentage": 10},
		"cost": {"total_cost_usd": 0.5},
		"rate_limits": {}
	}`
	git := &mockGitProvider{branch: "main"}
	var buf bytes.Buffer
	r := NewRenderer(&buf, git)

	require.NoError(t, r.GenerateFrom(t.Context(), strings.NewReader(json)))

	out := buf.String()
	assert.Contains(t, out, "claude-3")
	assert.Contains(t, out, "/tmp")
}

func TestRenderer_RenderFromReader_InvalidJSON(t *testing.T) {
	git := &mockGitProvider{}
	r := NewRenderer(&bytes.Buffer{}, git)
	err := r.GenerateFrom(t.Context(), strings.NewReader("{not json}"))
	assert.Error(t, err)
}
