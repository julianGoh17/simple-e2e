package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/internal"
	"github.com/julianGoh17/simple-e2e/framework/util"
	"github.com/stretchr/testify/assert"
)

func TestRunCmdValues(t *testing.T) {
	rootCmd := NewRootCmd()
	runCmd := NewRunCmd()
	initRunCmd(rootCmd, runCmd)

	assert.Equal(t, "run", runCmd.Use)
	assert.Equal(t, "Run through a set of or all the stages in a test", runCmd.Short)
	assert.Equal(t, "Run all the steps in a specified test or just a specific set of stages from that test.", runCmd.Long)
}

func TestRunCmdFailsWhenNoArgumentsPassedIn(t *testing.T) {
	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)

	// Unsure why but need to capture output through setting output. Have tried doing it other way, but it won't capture output
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"run"})
	assert.Error(t, rootCmd.Execute())

	out, err := ioutil.ReadAll(b)
	stringedOutput := string(out)

	assert.NoError(t, err)
	assert.Contains(t, stringedOutput, "Error: required flag(s) \"test\" not set")
}

func TestRunCmdFailsWhenCanNotFindFile(t *testing.T) {
	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)

	actualTestFilesDir := os.Getenv(util.TestDirEnv)
	os.Setenv(util.TestDirEnv, "")

	// Unsure why but need to capture output through setting output. Have tried doing it other way, but it won't capture output
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"run", "-t", "non-existent-test"})
	assert.Error(t, rootCmd.Execute())

	out, err := ioutil.ReadAll(b)
	stringedOutput := string(out)

	assert.NoError(t, err)
	assert.Contains(t, stringedOutput, "Error: unable to read file: open /non-existent-test.yaml")
	os.Setenv(util.TestDirEnv, actualTestFilesDir)

}

func TestRunCmdFailsWhenCanNotInitializeDockerHandler(t *testing.T) {
	os.Setenv(dockerHostEnv, invalidHost)
	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)

	rootCmd.SetArgs([]string{"run", "-t", "test"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Equal(t, invalidHostError, err.Error())
	os.Unsetenv(dockerHostEnv)
}

func TestRunCmdPassWhenCanFindValidTestFile(t *testing.T) {
	internal.SetTestFilesRoot()
	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)
	read, written, rescue := beginCaptureOfTerminalOutput()

	rootCmd.SetArgs([]string{"run", "-t", "test"})
	assert.NoError(t, rootCmd.Execute())

	output := endCaptureOfTerminalOutput(read, written, rescue)

	assert.Contains(t, output, "Hello there Julian!")
	assert.Contains(t, output, "Hello there Coachella!")
	assert.Contains(t, output, "Hello there Eugene!")
	assert.Contains(t, output, "Hello there Boy!")
}

func TestRunCmdPassWhenCanFindValidTestFileAndRunningFewStages(t *testing.T) {
	internal.SetTestFilesRoot()
	rootCmd := NewRootCmd()
	InitRootCmd(rootCmd)

	read, written, rescue := beginCaptureOfTerminalOutput()

	rootCmd.SetArgs([]string{"run", "-t", "test", "-s", "stage1"})
	assert.NoError(t, rootCmd.Execute())

	output := endCaptureOfTerminalOutput(read, written, rescue)

	assert.Contains(t, output, "Hello there Julian!")
	assert.Contains(t, output, "Hello there Coachella!")
	assert.NotContains(t, output, "Hello there Eugene!")
	assert.NotContains(t, output, "Hello there Boy!")
}

func beginCaptureOfTerminalOutput() (*os.File, *os.File, *os.File) {
	rescueStdout := os.Stdout
	read, written, _ := os.Pipe()
	os.Stdout = written

	return read, written, rescueStdout
}

func endCaptureOfTerminalOutput(read, written, rescueStdout *os.File) string {
	written.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = rescueStdout
	return string(out)
}
