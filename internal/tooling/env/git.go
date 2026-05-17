package env

import (
	"context"
	"os/exec"
	"strings"
)

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
	isGitDirty := checkGitDirty(ctx, dir)
	return branch, isGitDirty, nil
}

func checkGitDirty(ctx context.Context, repositoryDir string) bool {
	gitCmds := []*exec.Cmd{
		exec.CommandContext(ctx, "git", "-C", repositoryDir, "diff", "--quiet"),
		exec.CommandContext(ctx, "git", "-C", repositoryDir, "diff", "--cached", "--quiet"),
	}
	for _, cmd := range gitCmds {
		if err := cmd.Run(); err != nil {
			return true
		}
	}
	return false
}
