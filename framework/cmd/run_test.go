package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

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
	runCmd := NewRunCmd()
	initRunCmd(rootCmd, runCmd)

	// Unsure why but need to capture output through setting output. Have tried doing it other way, but it won't capture output
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"run"})
	rootCmd.Execute()

	out, err := ioutil.ReadAll(b)
	stringedOutput := string(out)

	assert.NoError(t, err)
	assert.Contains(t, stringedOutput, "Error: required flag(s) \"test\" not set")
}

func TestRunCmdFailsWhenCanNotFindFile(t *testing.T) {
	rootCmd := NewRootCmd()
	runCmd := NewRunCmd()
	initRunCmd(rootCmd, runCmd)

	// Unsure why but need to capture output through setting output. Have tried doing it other way, but it won't capture output
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"run", "-t", "test"})
	rootCmd.Execute()

	out, err := ioutil.ReadAll(b)
	stringedOutput := string(out)

	assert.NoError(t, err)
	assert.Contains(t, stringedOutput, "Error: unable to read file: open /test.yaml")
}

func TestRunCmdPassWhenCanFindValidTestFile(t *testing.T) {
	rootCmd := NewRootCmd()
	runCmd := NewRunCmd()
	initRunCmd(rootCmd, runCmd)
	testFileRoot := GetTestFilesRoot()
	os.Setenv(util.TestDirEnv, testFileRoot)
	read, written, rescue := beginCaptureOfTerminalOutput()

	rootCmd.SetArgs([]string{"run", "-t", "test"})
	rootCmd.Execute()

	output := endCaptureOfTerminalOutput(read, written, rescue)

	assert.Contains(t, output, "Hello there Julian!")
	assert.Contains(t, output, "Hello there Coachella!")
	assert.Contains(t, output, "Hello there Eugene!")
	assert.Contains(t, output, "Hello there Boy!")
	os.Unsetenv(util.TestDirEnv)
}

func TestRunCmdPassWhenCanFindValidTestFileAndRunningFewStages(t *testing.T) {
	rootCmd := NewRootCmd()
	runCmd := NewRunCmd()
	initRunCmd(rootCmd, runCmd)
	testFileRoot := GetTestFilesRoot()
	os.Setenv(util.TestDirEnv, testFileRoot)
	read, written, rescue := beginCaptureOfTerminalOutput()

	rootCmd.SetArgs([]string{"run", "-t", "test", "-s", "stage1"})
	rootCmd.Execute()

	output := endCaptureOfTerminalOutput(read, written, rescue)

	assert.Contains(t, output, "Hello there Julian!")
	assert.Contains(t, output, "Hello there Coachella!")
	assert.NotContains(t, output, "Hello there Eugene!")
	assert.NotContains(t, output, "Hello there Boy!")
}

func GetTestFilesRoot() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return fmt.Sprintf("%s/../tests", filepath.Dir(d))
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
