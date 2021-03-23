package main

import (
	"context"

	"github.com/4armed/killager/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Generate runs the generate command....
func Generate(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "generate",
		TraverseChildren: true,
		Short:            "Generate kubeconfig file with serviceAccount tokens found on node",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infof("processing secrets on node %s", c.Node)

			// Parse the kubeconfig file specified
			config, err := clientcmd.BuildConfigFromFlags("", c.KubeConfigFile)
			if err != nil {
				return err
			}

			// Create clientset
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}

			// Loop through the pods and get the secret volume mounts
			pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return err
			}
			logrus.Debugf("There are %d pods in the cluster\n", len(pods.Items))

			return nil
		},
	}

	cmd.Flags().StringVarP(&c.KubeConfigFile, "kubeconfig", "k", "kubeconfig.yaml", "The kubeconfig file to read cluster config from")
	cmd.Flags().StringVar(&c.Node, "node", "", "Node to process secrets for")
	cmd.MarkFlagRequired("node")

	return cmd
}
