/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package git

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/callsign/gitops/internal/testutil"
)

func Test_Checkout(t *testing.T) {
	tests := []struct {
		name          string
		projectURL    string
		environment   string
		expectedError error
		expectedFile  string
	}{{
		name:          "should return an error on invalid project url",
		projectURL:    "invalid",
		expectedError: fmt.Errorf("Cannot git clone*"),
	}, {
		name:          "should return an error on invalid environment",
		projectURL:    "https://github.com/githubtraining/hellogitworld.git",
		environment:   "invalid",
		expectedError: fmt.Errorf("Cannot git checkout*"),
	}, {
		name:         "should checkout the project",
		projectURL:   "https://github.com/githubtraining/hellogitworld.git",
		environment:  "feature_image",
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
	os.RemoveAll(directory)
}
