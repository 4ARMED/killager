package config

type Config struct {
	KubeConfigFile       string
	KubeConfigOutputFile string
	LogLevel             int
	Node                 string
}
