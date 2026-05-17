package render

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jgfranco17/delphi-cli/internal/tooling/git"
	"github.com/jgfranco17/delphi-cli/internal/tooling/model"
)

// Renderer writes a formatted statusline to an io.Writer.
type Renderer struct {
	out io.Writer
	git git.Provider
}

func NewRenderer(out io.Writer, g git.Provider) *Renderer {
	return &Renderer{out: out, git: g}
}

// New creates a Renderer with the real git provider writing to out.
func New(out io.Writer) *Renderer {
	return NewRenderer(out, git.NewExecProvider())
}

// GenerateFrom decodes JSON from r and renders the statusline.
func (r *Renderer) GenerateFrom(ctx context.Context, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed reading input: %w", err)
	}
	var input model.AgentInput
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("failed parsing input: %w", err)
	}
	return r.Render(ctx, &input)
}

// Render writes the three-line statusline for the given input.
func (r *Renderer) Render(ctx context.Context, input *model.AgentInput) error {
	pct := int(input.ContextWindow.UsedPercentage)
	bar := input.ContextWindow.Render()
	costFmt := input.Cost.Format()
	limits := input.RateLimits.Format()
	gitInfo := r.formatGitInfo(ctx, input.Workspace.CurrentDir)

	fmt.Fprintf(
		r.out, "%s %s %s %s\n",
		model.ColorDim.Sprint("Using"),
		model.ColorBoldCyan.Sprint(input.Model.DisplayName),
		model.ColorDim.Sprint("in"),
		model.ColorYellow.Sprint(input.Workspace.CurrentDir),
	)
	fmt.Fprintf(
		r.out, "%s %s %s %s %s %s %s\n",
		model.ColorDim.Sprint("Usage:"),
		input.ContextWindow.Color().Sprint(bar),
		model.ColorBold.Sprintf("%d%%", pct),
		model.ColorDim.Sprint("|"),
		model.ColorGreen.Sprintf("~%s equiv", costFmt),
		model.ColorDim.Sprint("|"),
		model.ColorMagenta.Sprint(limits),
	)
	fmt.Fprintf(
		r.out, "%s %s\n",
		model.ColorDim.Sprint("Git:"),
		gitInfo,
	)
	return nil
}

func (r *Renderer) formatGitInfo(ctx context.Context, dir string) string {
	branch, dirty, err := r.git.BranchStatus(ctx, dir)
	if err != nil || branch == "" {
		return "none"
	}
	branchStr := model.ColorBoldCyan.Sprint(branch)
	if dirty {
		return fmt.Sprintf("%s %s", branchStr, model.ColorRed.Sprint("✗ dirty"))
	}
	return fmt.Sprintf("%s %s", branchStr, model.ColorGreen.Sprint("✓ clean"))
}
