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

package environment

import (
	"fmt"
	"path"
	"testing"

	"github.com/callsign/gitops/internal/testutil"
)

func Test_Get(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{{
		name:     "should return the environment",
		expected: "prod",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Get()
			if err != nil {
				t.Fatalf("\nTest failure: %v", err)
			}
			if actual != test.expected {
				t.Fatalf("\nUnexpected result:\nExpected: %v\nGot: %v", test.expected, actual)
			}
		})
	}
}

func Test_serviceBranchToEnvironment(t *testing.T) {
	tests := []struct {
		name                  string
		serviceBranch         string
		customDeploymentsYaml string
		expectedError         error
		expectedResult        string
	}{{
		name:           "should support standard environments",
		serviceBranch:  "develop",
		expectedResult: "dev",
	}, {
		name:          "should return an error on failure to determine the deployment environment",
		serviceBranch: "unknown",
		expectedError: fmt.Errorf("Cannot find deployment environment for service branch unknown"),
	}, {
		name:                  "should support custom deployments",
		serviceBranch:         "feature/foo",
		customDeploymentsYaml: "valid-custom-deployments.yaml",
		expectedResult:        "bar",
	}, {
		name:                  "should not allow custom deployments to deploy to standard environments",
		serviceBranch:         "feature/foo",
		customDeploymentsYaml: "forbidden-custom-deployments.yaml",
		expectedError:         fmt.Errorf("Deployment of service branch feature/foo to deployment environment prod forbidden"),
	}, {
		name: "should return an error on unreadable custom deployments yaml",
		customDeploymentsYaml: ".",
		expectedError:         fmt.Errorf("Cannot read *"),
	}, {
		name:                  "should return an error on invalid custom deployments yaml",
		serviceBranch:         "feature/foo",
		customDeploymentsYaml: "invalid-custom-deployments.yaml",
		expectedError:         fmt.Errorf("Cannot unmarshal *"),
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var customDeploymentsYamlPath string
			if test.customDeploymentsYaml != "" {
				customDeploymentsYamlPath = path.Join(testutil.Data("environment/get"), test.customDeploymentsYaml)
			}
			actualResult, actualError := serviceBranchToEnvironment(test.serviceBranch, customDeploymentsYamlPath)
			testutil.VerifyError(test.expectedError, actualError, t)
			if test.expectedResult != actualResult {
				t.Fatalf("\nUnexpected result:\nExpected: %v\nGot: %v", test.expectedResult, actualResult)
			}
		})
	}
}
