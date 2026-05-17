package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecProvider_NotARepo(t *testing.T) {
	dir := t.TempDir()
	p := NewExecProvider()
	branch, _, err := p.BranchStatus(t.Context(), dir)
	assert.Error(t, err)
	assert.Empty(t, branch)
}

func TestExecProvider_CleanRepo(t *testing.T) {
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

	p := NewExecProvider()
	branch, dirty, err := p.BranchStatus(t.Context(), dir)
	require.NoError(t, err)
	assert.Equal(t, "main", branch)
	assert.False(t, dirty)
}

func TestExecProvider_DirtyRepo(t *testing.T) {
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

	require.NoError(t, os.WriteFile(readmeFile, []byte("changed"), 0644))

	p := NewExecProvider()
	branch, dirty, err := p.BranchStatus(t.Context(), dir)
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
