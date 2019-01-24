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

package service

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/callsign/gitops/internal/directory"
	"github.com/callsign/gitops/internal/testutil"
)

func Test_Update(t *testing.T) {
	tests := []struct {
		name                     string
		serviceName              string
		serviceVersion           string
		input                    string
		expectedError            error
		expectedRequirementsYaml string
		expectedRequirementsLock string
	}{{
		name:          "should return an error if the service name is missing",
		expectedError: fmt.Errorf("Missing service name"),
	}, {
		name:          "should return an error if the service version is missing",
		serviceName:   "service",
		expectedError: fmt.Errorf("Missing service version"),
	}, {
		name:           "should return an error if the charts directory is missing",
		serviceName:    "service",
		serviceVersion: "2.1.1805200912+abcde",
		input:          "missing-charts-directory",
		expectedError:  fmt.Errorf("Cannot find chart: Missing charts directory"),
	}, {
		name:           "should return an error if the chart is missing",
		serviceName:    "service",
		serviceVersion: "2.1.1805200912+abcde",
		input:          "missing-chart",
		expectedError:  fmt.Errorf("Cannot find chart: Missing chart"),
	}, {
		name:           "should add a dependency if the file is missing",
		serviceName:    "service",
		serviceVersion: "2.1.1805200912+abcde",
		input:          "missing-requirements",
		expectedRequirementsYaml: "expected-requirements.yaml",
	}, {
		name:           "should add a dependency if the file is presant but the service is missing",
		serviceName:    "service",
		serviceVersion: "2.1.1805200912+abcde",
		input:          "missing-service",
		expectedRequirementsYaml: "expected-requirements.yaml",
	}, {
		name:           "should update the version and the requirements.lock if the file and the service are present",
		serviceName:    "service",
		serviceVersion: "2.1.1805200912+abcde",
		input:          "version-updated-lock-file-updated",
		expectedRequirementsYaml: "expected-requirements.yaml",
		expectedRequirementsLock: "expected-requirements.lock",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testData := path.Join(testutil.Data("service/update"), test.input)
			generatedTestData := path.Join(testutil.GeneratedData("service/update", true), test.input)
			os.RemoveAll(generatedTestData)
			if err := directory.Copy(testData, generatedTestData); err != nil {
				t.Fatalf("\nCannot copy test data: %v", err)
			}

			projectName := path.Join(testutil.GeneratedData("service/update", false), test.input)
			companyHelmRepository = "file://../../service/charts/service"
			err := Update(projectName, test.serviceName, test.serviceVersion)

			testutil.VerifyError(test.expectedError, err, t)

			expected := testutil.ReadFile(path.Join(testData, test.expectedRequirementsYaml), t, err)
			actual := testutil.ReadFile(path.Join(generatedTestData, "/charts/application/requirements.yaml"), t, err)
			if err == nil && actual != expected {
				t.Fatalf("\nUnexpected requirements.yaml:\nExpected: %v\nGot: %v", expected, actual)
			}

			if test.expectedRequirementsLock != "" {
				expected = testutil.ReadFile(path.Join(testData, test.expectedRequirementsLock), t, err)
				actual = testutil.ReadFile(path.Join(generatedTestData, "/charts/application/requirements.lock"), t, err)
				if err == nil && !strings.HasPrefix(actual, expected) {
					t.Fatalf("\nUnexpected requirements.lock:\nExpected: %v\nGot: %v", expected, actual)
				}
			}
		})
	}
}
