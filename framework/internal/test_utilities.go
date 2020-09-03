package internal

/*
 * THIS PACKAGE IS MEANT FOR TEST UTILITY FUNCTIONS ONLY
 */

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/julianGoh17/simple-e2e/framework/util"
)

// SetDockerfilesRoot will set 'DOCKERFILES_DIR' env to the path to the 'project's root/Dockerfiles' if it's not already set
func SetDockerfilesRoot() {
	// If not in container, set as the path to the 'project's root/Dockerfiles'
	if os.Getenv(util.DockerfileDirEnv) == "" {
		_, b, _, _ := runtime.Caller(0)
		d := path.Join(path.Dir(b))
		os.Setenv(util.DockerfileDirEnv, fmt.Sprintf("%s/../Dockerfiles", filepath.Dir(d)))
	}
}

// SetTestFilesRoot will set 'TEST_DIR' env to the path to the 'project's root/tests' if it's not already set
func SetTestFilesRoot() {
	// If not in container, set as the path to the 'project's root/tests'
	if os.Getenv(util.TestDirEnv) == "" {
		_, b, _, _ := runtime.Caller(0)
		d := path.Join(path.Dir(b))
		os.Setenv(util.TestDirEnv, fmt.Sprintf("%s/../tests", filepath.Dir(d)))
	}
}

// TestCoverageReaches85Percent will ensure that test coverage passes 85%
func TestCoverageReaches85Percent(m *testing.M) {
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
