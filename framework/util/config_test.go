package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrDefaultValues(t *testing.T) {
	values := []struct {
		setValue      string
		defaultValue  string
		expectedValue string
	}{
		{
			"",
			"default",
			"default",
		},
		{
			"set",
			"",
			"set",
		},
	}

	config := NewConfig()
	key := "REALLY_RANDOM_ENVIRONMENTAL_VARIABLE"
	for _, value := range values {
		os.Setenv(key, value.setValue)
		config.setDefault(key, value.defaultValue)
		assert.Equal(t, value.expectedValue, config.GetOrDefault(key))
		os.Unsetenv(key)
	}

}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.85 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
