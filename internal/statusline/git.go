package statusline

import (
	"os/exec"
	"strings"
)

// GitProvider abstracts git operations for testability.
type GitProvider interface {
	BranchStatus(dir string) (branch string, dirty bool, err error)
}

// ExecGitProvider is the real implementation that shells out to git.
type ExecGitProvider struct{}

func NewExecGitProvider() *ExecGitProvider {
	return &ExecGitProvider{}
}

// BranchStatus returns the current branch name and whether the working tree is dirty.
// Returns an error (and empty branch) when dir is not a git repository.
func (g *ExecGitProvider) BranchStatus(dir string) (string, bool, error) {
	out, err := exec.Command("git", "-C", dir, "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return "", false, err
	}
	branch := strings.TrimSpace(string(out))

	diffErr := exec.Command("git", "-C", dir, "diff", "--quiet").Run()
	cachedErr := exec.Command("git", "-C", dir, "diff", "--cached", "--quiet").Run()
	dirty := diffErr != nil || cachedErr != nil

	return branch, dirty, nil
}
