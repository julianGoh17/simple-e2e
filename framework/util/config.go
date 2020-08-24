package util

import "os"

// GlobalConfig should be used when the framework is trying to get an environment variable from the container
type GlobalConfig struct{}

// GetOrDefault will get a value from the environment and when it can't find the value from the environment it will return the default value
func (config *GlobalConfig) GetOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
