package main

import (
	"fmt"
	"os"

	"github.com/4armed/killager/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Configuration struct
var c *config.Config

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

	var verboseLogging bool

	rootCmd.PersistentFlags().BoolVarP(&verboseLogging, "verbose", "v", false, "Output debug statements")

	if verboseLogging {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	rootCmd.AddCommand(Generate(c))
}
