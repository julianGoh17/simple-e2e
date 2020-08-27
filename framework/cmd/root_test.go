package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	rootCmd := NewRootCmd()

	assert.Equal(t, rootCmd.Use, "Simple-E2E")
	assert.Equal(t, rootCmd.Short, "A modular and configurable testing infrastructure")
	assert.Equal(t, rootCmd.Long, `Simple-E2E is a testing library aimed at making more modular and easier. 
		This application allows users to break down tests into stages and steps to 
		run a set of stages or an entire test. Furthermore, Simple-E2E provides a
		framework to easily create new tests from exisiting steps.`)
}

func TestInitRootCmd(t *testing.T) {
	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)

	assert.Equal(t, 2, len(rootCmd.Commands()))
	assert.Equal(t, "run", rootCmd.Commands()[0].Use)
	assert.Equal(t, "version", rootCmd.Commands()[1].Use)
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
