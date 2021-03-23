package main

import (
	"github.com/4armed/killager/pkg/config"
	"github.com/spf13/cobra"
)

// Generate runs the generate command....
func Generate(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "generate",
		TraverseChildren: true,
		Short:            "Generate kubeconfig file with serviceAccount tokens found on node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVarP(&c.KubeConfigFile, "kubeconfig", "k", "kubeconfig.yaml", "The kubeconfig file to read cluster config from")
	cmd.Flags().StringVar(&c.Node, "node", "", "Node to process secrets for")

	return cmd
}
