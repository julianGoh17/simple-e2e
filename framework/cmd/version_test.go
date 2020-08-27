package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	rootCmd := NewRootCmd()
	versionCmd := NewVersionCmd()
	initVersionCmd(rootCmd, versionCmd)
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"version"})
	rootCmd.Execute()
	out, err := ioutil.ReadAll(b)
	stringedOutput := string(out)

	assert.NoError(t, err)
	assert.Equal(t, stringedOutput, "Simple-E2E binary version: v0.1")
}
