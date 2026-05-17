package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jgfranco17/delphi-cli/internal/tooling/render"
	"github.com/spf13/cobra"
)

const (
	noColorFlagName   = "disable-colors"
	showAllFlagName   = "all"
	showCwdFlagName   = "cwd"
	showGitFlagName   = "git"
	showUsageFlagName = "stats"
)

func init() {
	color.NoColor = false
}

func newStatuslineCmd() *cobra.Command {
	var noColor bool
	var showAll bool
	var showCwd bool
	var showGit bool
	var showUsage bool

	cmd := &cobra.Command{
		Use:   "statusline",
		Short: "Render Claude Code statusline from JSON piped in stdin",
		Long: `Render Claude Code statusline from JSON piped in stdin.

The statusline command reads a JSON blob from stdin representing the agent's state and renders a formatted statusline to stdout.
This is intended for use by shell integrations to display the current model, workspace, git status, and usage stats.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if noColor {
				color.NoColor = true
			}
			opts := render.Options{
				OutputWriter:   cmd.OutOrStdout(),
				WithCurrentDir: showAll || showCwd,
				WithGitInfo:    showAll || showGit,
				WithUsageStats: showAll || showUsage,
			}
			status := render.New(opts)
			if err := status.GenerateFrom(cmd.Context(), os.Stdin); err != nil {
				return fmt.Errorf("failed to render status line: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&noColor, noColorFlagName, false, "Disable ANSI color output formatting")
	cmd.Flags().BoolVar(&showAll, showAllFlagName, false, "Show all statusline sections")
	cmd.Flags().BoolVar(&showCwd, showCwdFlagName, false, "Show current working directory")
	cmd.Flags().BoolVar(&showGit, showGitFlagName, false, "Show git information")
	cmd.Flags().BoolVar(&showUsage, showUsageFlagName, false, "Show usage statistics")
	return cmd
}
