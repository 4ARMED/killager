package config

type Config struct {
	KubeConfigFile       string
	KubeConfigOutputFile string
	VerboseLogging       bool
	Node                 string
}
