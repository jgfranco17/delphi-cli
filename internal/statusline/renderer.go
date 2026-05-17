package statusline

import (
	"encoding/json"
	"fmt"
	"io"
)

// Renderer writes a formatted statusline to an io.Writer.
type Renderer struct {
	out io.Writer
	git GitProvider
}

func NewRenderer(out io.Writer, git GitProvider) *Renderer {
	return &Renderer{out: out, git: git}
}

// RenderFromReader decodes JSON from r and renders the statusline.
func (r *Renderer) RenderFromReader(reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}
	var input Input
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("parsing input: %w", err)
	}
	return r.Render(&input)
}

// Render writes the three-line statusline for the given input.
func (r *Renderer) Render(input *Input) error {
	pct := int(input.ContextWindow.UsedPercentage)
	bar := input.ContextWindow.Render()
	costFmt := input.Cost.Format()
	limits := input.RateLimits.Format()
	gitInfo := r.formatGitInfo(input.Workspace.CurrentDir)

	fmt.Fprintf(r.out, "%s %s %s %s\n",
		colorDim.Sprint("Using"),
		colorBoldCyan.Sprint(input.Model.DisplayName),
		colorDim.Sprint("in"),
		colorYellow.Sprint(input.Workspace.CurrentDir))
	fmt.Fprintf(r.out, "%s %s %s %s %s %s %s\n",
		colorDim.Sprint("Usage:"),
		input.ContextWindow.Color().Sprint(bar),
		colorBold.Sprintf("%d%%", pct),
		colorDim.Sprint("|"),
		colorGreen.Sprintf("~%s equiv", costFmt),
		colorDim.Sprint("|"),
		colorMagenta.Sprint(limits))
	fmt.Fprintf(r.out, "%s %s\n",
		colorDim.Sprint("Git:"),
		gitInfo)

	return nil
}

func (r *Renderer) formatGitInfo(dir string) string {
	branch, dirty, err := r.git.BranchStatus(dir)
	if err != nil || branch == "" {
		return "none"
	}
	branchStr := colorBoldCyan.Sprint(branch)
	if dirty {
		return fmt.Sprintf("%s %s", branchStr, colorRed.Sprint("✗ dirty"))
	}
	return fmt.Sprintf("%s %s", branchStr, colorGreen.Sprint("✓ clean"))
}
