package main

import (
	"fmt"
	"os"

	"github.com/4armed/killager/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// things we need
var c *config.Config
var verboseLogging bool
var silentLogging bool

var rootCmd = &cobra.Command{
	Version:       config.GitVersion,
	Use:           config.Executable,
	Short:         "Read service account secrets from target Kubernetes node and write a kubeconfig to use it",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	c = &config.Config{}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := initLogs(verboseLogging); err != nil {
			return err
		}
		return nil
	}
	rootCmd.PersistentFlags().BoolVarP(&silentLogging, "quiet", "q", false, "Suppress all log statements")
	rootCmd.PersistentFlags().BoolVarP(&verboseLogging, "verbose", "v", false, "Output debug statements")

	rootCmd.AddCommand(Generate(c))
}

func initLogs(verbose bool) error {
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else if silentLogging {
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	return nil

}
