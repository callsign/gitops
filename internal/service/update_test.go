// Copyright © 2018 Callsign. All rights reserved.

package service

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"callsign.com/gitops/internal/directory"
	"callsign.com/gitops/internal/testutil"
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
		name:           "should return an error if multiple charts exist",
		serviceName:    "service",
		serviceVersion: "2.1.1805200912+abcde",
		input:          "multiple-charts",
		expectedError:  fmt.Errorf("Cannot find chart: Multiple charts"),
	}, {
		name:           "should return an error if the file is missing",
		serviceName:    "service",
		serviceVersion: "2.1.1805200912+abcde",
		input:          "missing-requirements",
		expectedError:  fmt.Errorf("Cannot read requirements.yaml: Missing requirements.yaml"),
	}, {
		name:           "should return an error if the file is presant but the service is missing",
		serviceName:    "service",
		serviceVersion: "2.1.1805200912+abcde",
		input:          "missing-service",
		expectedError:  fmt.Errorf("Cannot update requirements.yaml: Missing service entry"),
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
			err := Update(projectName, test.serviceName, test.serviceVersion)

			testutil.VerifyError(test.expectedError, err, t)

			expected := testutil.ReadFile(path.Join(testData, test.expectedRequirementsYaml), t, err)
			actual := testutil.ReadFile(path.Join(generatedTestData, "/charts/application/requirements.yaml"), t, err)
			if err == nil && actual != expected {
				t.Fatalf("\nUnexpected requirements.yaml:\nExpected: %v\nGot: %v", expected, actual)
			}

			expected = testutil.ReadFile(path.Join(testData, test.expectedRequirementsLock), t, err)
			actual = testutil.ReadFile(path.Join(generatedTestData, "/charts/application/requirements.lock"), t, err)
			if err == nil && !strings.HasPrefix(actual, expected) {
				t.Fatalf("\nUnexpected requirements.lock:\nExpected: %v\nGot: %v", expected, actual)
			}
		})
	}
}
