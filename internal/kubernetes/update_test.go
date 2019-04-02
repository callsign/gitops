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

package kubernetes

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/callsign/gitops/internal/testutil"
)

func Test_Update(t *testing.T) {
	tests := []struct {
		name              string
		projectName       string
		serviceName       string
		environment       string
		expectedError     error
		expectedResources string
	}{{
		name:        "should produce kubernetes resources on success",
		projectName: "success",
		serviceName: "test",
		environment: "dev",
	}, {
		name:          "should return an error if the project name is missing",
		expectedError: fmt.Errorf("Missing project name"),
	}, {
		name:          "should return an error if the service name is missing",
		projectName:   "missing-service-name",
		expectedError: fmt.Errorf("Missing service name"),
	}, {
		name:          "should return an error if the environment is missing",
		projectName:   "missing-environment",
		serviceName:   "test",
		expectedError: fmt.Errorf("Missing environment"),
	}, {
		name:          "should return an error on helm template execution error",
		projectName:   "helm-template-execution-error",
		serviceName:   "test",
		environment:   "dev",
		expectedError: fmt.Errorf("Cannot execute helm template: *"),
	}, {
		name:          "should return an error on missing packaged chart directory",
		projectName:   "missing-packaged-chart-directory",
		serviceName:   "test",
		environment:   "dev",
		expectedError: fmt.Errorf("Cannot execute helm template: Cannot read build/packages/helm directory"),
	}, {
		name:          "should return an error on missing packaged chart",
		projectName:   "missing-packaged-chart",
		serviceName:   "test",
		environment:   "dev",
		expectedError: fmt.Errorf("Cannot execute helm template: Missing packaged chart in build/packages/helm"),
	}, {
		name:          "should return an error if the project directory does not exist",
		projectName:   "project-directory-does-not-exist",
		serviceName:   "test",
		environment:   "dev",
		expectedError: fmt.Errorf("Cannot copy Kubernetes resources: Project directory does not exist (*"),
	}, {
		name:          "should return an error if the project directory is a file",
		projectName:   "project-directory-is-a-file",
		serviceName:   "test",
		environment:   "dev",
		expectedError: fmt.Errorf("Cannot copy Kubernetes resources: Project directory is a file (*"),
	}, {
		name:          "should return an error on resources copy error",
		projectName:   "resources-copy-error",
		serviceName:   "test",
		environment:   "dev",
		expectedError: fmt.Errorf("Cannot copy Kubernetes resources: *"),
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			previousCurrentDirectory, _ := os.Getwd()

			projectName := ""
			if test.projectName != "" {
				projectName = path.Join(testutil.GeneratedData("kubernetes/update", true), test.projectName)
				os.RemoveAll(projectName)
				if test.projectName == "project-directory-is-a-file" {
					emptyFile, _ := os.Create(projectName)
					emptyFile.Close()
				} else if test.projectName == "resources-copy-error" {
					os.MkdirAll(projectName, 0555)
				} else if test.projectName != "project-directory-does-not-exist" {
					os.MkdirAll(projectName, 0755)
				}
				_ = os.Chdir(path.Join(testutil.Data("kubernetes/update"), test.projectName))
			}

			err := Update(projectName, test.serviceName, test.environment)

			_ = os.Chdir(previousCurrentDirectory)

			testutil.VerifyError(test.expectedError, err, t)

			if err == nil {
				expectedFile := path.Join(projectName, "kubernetes", "test", "deployment.yaml")
				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
					t.Fatalf("\nExpected file %s does not exist", expectedFile)
				}
			}
		})
	}
}