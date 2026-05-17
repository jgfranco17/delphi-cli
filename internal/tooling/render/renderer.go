package render

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jgfranco17/delphi-cli/internal/tooling/git"
	"github.com/jgfranco17/delphi-cli/internal/tooling/model"
)

// StatusLine writes a formatted statusline to a writer.
type StatusLine struct {
	out     io.Writer
	git     git.Provider
	options Options
}

type Options struct {
	OutputWriter   io.Writer
	WithCurrentDir bool
	WithGitInfo    bool
	WithUsageStats bool
}

func newRenderer(out io.Writer, g git.Provider, opts Options) *StatusLine {
	return &StatusLine{out: out, git: g, options: opts}
}

// New creates a StatusLine with the real git provider writing to out.
func New(opts Options) *StatusLine {
	out := opts.OutputWriter
	if out == nil {
		out = os.Stdout
	}
	return newRenderer(out, git.NewExecProvider(), opts)
}

// GenerateFrom decodes JSON from r and renders the statusline.
func (r *StatusLine) GenerateFrom(ctx context.Context, reader io.Reader) error {
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

// Render writes the statusline for the given input, gated by Options.
func (r *StatusLine) Render(ctx context.Context, input *model.AgentInput) error {
	if r.options.WithCurrentDir {
		fmt.Fprintf(
			r.out, "%s %s %s %s\n",
			model.ColorDim.Sprint("Using"),
			model.ColorBoldCyan.Sprint(input.Model.DisplayName),
			model.ColorDim.Sprint("in"),
			model.ColorYellow.Sprint(input.Workspace.CurrentDir),
		)
	} else {
		fmt.Fprintf(
			r.out, "%s %s\n",
			model.ColorDim.Sprint("Using"),
			model.ColorBoldCyan.Sprint(input.Model.DisplayName),
		)
	}
	if r.options.WithUsageStats {
		pct := int(input.ContextWindow.UsedPercentage)
		bar := input.ContextWindow.Render()
		costFmt := input.Cost.Format()
		limits := input.RateLimits.Format()
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
	}
	if r.options.WithGitInfo {
		gitInfo := r.formatGitInfo(ctx, input.Workspace.CurrentDir)
		fmt.Fprintf(
			r.out, "%s %s\n",
			model.ColorDim.Sprint("Git:"),
			gitInfo,
		)
	}
	return nil
}

func (r *StatusLine) formatGitInfo(ctx context.Context, dir string) string {
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
