package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jgfranco17/delphi-cli/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCLI() *cobra.Command {
	var verbosity int

	rootCmd := &cobra.Command{
		Use:   "delphi",
		Short: "Delphi CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var level logrus.Level
			switch verbosity {
			case 1:
				level = logrus.InfoLevel
			case 2:
				level = logrus.DebugLevel
			case 3:
				level = logrus.TraceLevel
			default:
				level = logrus.WarnLevel
			}

			logger := logging.New(cmd.ErrOrStderr(), level)
			ctx := logging.AddToContext(cmd.Context(), logger)

			ctx, cancel := context.WithCancel(ctx)
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				select {
				case sig := <-c:
					logger.WithField("signal", sig).Warn("Received signal, exiting")
					cancel()
				case <-ctx.Done():
					// Context was canceled, exit goroutine.
				}
			}()
			return nil
		},
	}

	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "Increase verbosity (up to -vvv)")

	rootCmd.AddCommand(newStatuslineCmd())
	return rootCmd
}

func Execute() {
	ctx := context.Background()
	cli := NewCLI()
	if err := cli.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
