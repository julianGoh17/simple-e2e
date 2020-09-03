package util

import "os"

// GlobalConfig should be used when the framework is trying to get an environment variable from the container
type GlobalConfig struct {
	defaults map[string]string
}

const (
	// TestDirEnv is the env var key for the root test file directory
	TestDirEnv = "TEST_DIR"
	// DockerfileDirEnv is the env var key for the root Dockerfile directory
	DockerfileDirEnv = "DOCKERFILE_DIR"
)

// NewConfig object returns the config object initialized with the default values
func NewConfig() *GlobalConfig {
	config := GlobalConfig{}
	initializeConfig(&config)
	return &config
}

// GetOrDefault will get a value from the environment and when it can't find the value from the environment it will return the default value
func (config *GlobalConfig) GetOrDefault(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return config.defaults[key]
}

func initializeConfig(config *GlobalConfig) {
	config.defaults = map[string]string{
		TestDirEnv:       "/home/e2e/tests",
		DockerfileDirEnv: "/home/e2e/Dockerfiles",
	}
}

// (TESTING ONLY) will override the given default with a new default
func (config *GlobalConfig) setDefault(key, value string) {
	config.defaults[key] = value
}
