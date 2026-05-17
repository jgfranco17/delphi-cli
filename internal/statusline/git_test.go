package statusline

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecGitProvider_NotARepo(t *testing.T) {
	dir := t.TempDir()
	p := NewExecGitProvider()
	branch, _, err := p.BranchStatus(dir)
	assert.Error(t, err)
	assert.Empty(t, branch)
}

func TestExecGitProvider_CleanRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	mustGit(t, dir, "init", "-b", "main")
	mustGit(t, dir, "config", "user.email", "test@example.com")
	mustGit(t, dir, "config", "user.name", "Test")

	// Need at least one commit for symbolic-ref to resolve.
	readmeFile := filepath.Join(dir, "README.md")
	require.NoError(t, os.WriteFile(readmeFile, []byte("hello"), 0644))
	mustGit(t, dir, "add", ".")
	mustGit(t, dir, "commit", "-m", "init")

	p := NewExecGitProvider()
	branch, dirty, err := p.BranchStatus(dir)
	require.NoError(t, err)
	assert.Equal(t, "main", branch)
	assert.False(t, dirty)
}

func TestExecGitProvider_DirtyRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	dir := t.TempDir()
	mustGit(t, dir, "init", "-b", "main")
	mustGit(t, dir, "config", "user.email", "test@example.com")
	mustGit(t, dir, "config", "user.name", "Test")

	readmeFile := filepath.Join(dir, "README.md")
	require.NoError(t, os.WriteFile(readmeFile, []byte("hello"), 0644))
	mustGit(t, dir, "add", ".")
	mustGit(t, dir, "commit", "-m", "init")

	// Modify a tracked file to make it dirty.
	require.NoError(t, os.WriteFile(readmeFile, []byte("changed"), 0644))

	p := NewExecGitProvider()
	branch, dirty, err := p.BranchStatus(dir)
	require.NoError(t, err)
	assert.Equal(t, "main", branch)
	assert.True(t, dirty)
}

func mustGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)
}
