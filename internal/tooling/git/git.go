package git

import (
	"context"
	"os/exec"
	"strings"
)

// Provider abstracts git operations for testability.
type Provider interface {
	BranchStatus(ctx context.Context, dir string) (branch string, dirty bool, err error)
}

// ExecProvider is the real implementation that shells out to git.
type ExecProvider struct{}

func NewExecProvider() *ExecProvider {
	return &ExecProvider{}
}

// BranchStatus returns the current branch name and whether the working tree is dirty.
// Returns an error (and empty branch) when dir is not a git repository.
func (g *ExecProvider) BranchStatus(ctx context.Context, dir string) (string, bool, error) {
	out, err := exec.CommandContext(ctx, "git", "-C", dir, "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return "", false, err
	}
	branch := strings.TrimSpace(string(out))

	diffErr := exec.CommandContext(ctx, "git", "-C", dir, "diff", "--quiet").Run()
	cachedErr := exec.CommandContext(ctx, "git", "-C", dir, "diff", "--cached", "--quiet").Run()
	dirty := diffErr != nil || cachedErr != nil

	return branch, dirty, nil
}
