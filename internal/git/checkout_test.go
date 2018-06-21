// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"callsign.com/gitops/internal/testutil"
)

func Test_Checkout(t *testing.T) {
	tests := []struct {
		name          string
		projectURL    string
		environment   string
		expectedError error
		expectedFile  string
	}{{
		name: "should return an error on invalid project url",
		projectURL: "invalid",
		expectedError: fmt.Errorf("Cannot clone *"),
	}, {
		name: "should return an error on invalid environment",
		projectURL: "https://github.com/githubtraining/hellogitworld.git",
		environment: "invalid",
		expectedError: fmt.Errorf("Cannot checkout *"),
	}, {
		name: "should checkout the project",
		projectURL: "https://github.com/githubtraining/hellogitworld.git",
		environment: "feature_image",
		expectedFile: "README.txt",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanup()

			projectName, actualError := Checkout(test.projectURL, test.environment)

			if actualError == nil {
				expectedFilePath, _ := filepath.Abs(path.Join(projectName, test.expectedFile))
				if _, err := os.Stat(expectedFilePath); err != nil {
					cleanup()
					t.Fatalf("Cannot find expected file (%s)", expectedFilePath)
				} else {
					cleanup()
				}
			} else {
				cleanup()
			}

			testutil.VerifyError(test.expectedError, actualError, t)
		})
	}
}

func cleanup() {
	directory, _ := filepath.Abs("hellogitworld")
	fmt.Println("directory:" + directory)
	os.RemoveAll(directory)
}
