package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/jgfranco17/delphi-cli/internal/tooling/render"
	"github.com/spf13/cobra"
)

func newStatuslineCmd() *cobra.Command {
	const forceColorFlagName = "force-color"
	var forceColor bool

	cmd := &cobra.Command{
		Use:   "statusline",
		Short: "Render Claude Code statusline from JSON piped on stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			forceColor, _ := cmd.Flags().GetBool(forceColorFlagName)
			if forceColor {
				color.NoColor = false
			}
			renderer := render.New(os.Stdout)
			return renderer.GenerateFrom(cmd.Context(), os.Stdin)
		},
	}

	cmd.Flags().BoolVar(&forceColor, forceColorFlagName, false, "Force ANSI color output formatting even when stdout is not a TTY")
	return cmd
}
