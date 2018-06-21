// Copyright © 2018 Callsign. All rights reserved.

package configuration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"callsign.com/gitops/internal/testutil"
)

func Test_Update(t *testing.T) {
	tests := []struct {
		name                      string
		configurationPath         string
		environment               string
		projectNameMissing        bool
		projectDirectoryMissing   bool
		projectDirectoryIsFile    bool
		serviceName               string
		input                     string
		expectedError             error
		expectedConfigurationFile string
	}{{
		name:          "should return an error if the configuration directory is missing",
		environment:   "staging",
		serviceName:   "service",
		expectedError: fmt.Errorf("Cannot copy service configuration: Configuration path missing"),
	}, {
		name:              "should return an error if the configuration directory does not exist",
		configurationPath: "configuration-directory-does-not-exist",
		environment:       "staging",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Source directory does not exist (*"),
	}, {
		name:              "should return an error if the configuration directory is a file",
		configurationPath: "configuration-directory-is-file",
		environment:       "staging",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Source directory does not exist (*"),
	}, {
		name:              "should return an error if the environment is missing",
		configurationPath: "environment-missing",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Environment missing"),
	}, {
		name:              "should return an error if the environment directory does not exist",
		configurationPath: "environment-directory-does-not-exist",
		environment:       "staging",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Source directory does not exist (*"),
	}, {
		name:              "should return an error if the environment directory is a file",
		configurationPath: "environment-directory-is-file",
		environment:       "staging",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Source directory is a file (*"),
	}, {
		name:               "should return an error if the projectName is missing",
		configurationPath:  "project-name-missing",
		environment:        "staging",
		projectNameMissing: true,
		serviceName:        "service",
		expectedError:      fmt.Errorf("Cannot copy service configuration: Project name missing"),
	}, {
		name:                    "should return an error if the project directory does not exist",
		configurationPath:       "project-directory-does-not-exist",
		environment:             "staging",
		projectDirectoryMissing: true,
		serviceName:             "service",
		expectedError:           fmt.Errorf("Cannot copy service configuration: Project directory does not exist (*"),
	}, {
		name:                    "should return an error if the project directory is a file",
		configurationPath:       "project-directory-is-file",
		environment:             "staging",
		projectDirectoryMissing: true,
		projectDirectoryIsFile:  true,
		serviceName:             "service",
		expectedError:           fmt.Errorf("Cannot copy service configuration: Project directory is a file (*"),
	}, {
		name:              "should return an error if the serviceName is missing",
		configurationPath: "service-name-missing",
		environment:       "staging",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Service name missing"),
	}, {
		name:              "should return an error on invalid yaml",
		configurationPath: "invalid-yaml",
		environment:       "staging",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot add service prefix to service configuration: Invalid YAML (*"),
	}, {
		name:                      "should copy the configuration and add a service prefix to the files",
		configurationPath:         "configuration-copied-service-prefix-added",
		environment:               "staging",
		serviceName:               "service",
		expectedConfigurationFile: "expected-configuration-file.yaml",
	}, {
		name:                      "should ignore non yaml files and directories",
		configurationPath:         "ignore-non-yaml-files-and-directories",
		environment:               "staging",
		serviceName:               "service",
		expectedConfigurationFile: "expected-configuration-file.yaml",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var configurationPath string
			var projectName string
			var testData string
			if test.configurationPath != "" {
				testData = path.Join(testutil.Data("configuration/update"), test.configurationPath)
				configurationPath = path.Join(testData, "environments")
				if !test.projectNameMissing {
					projectName = path.Join(testutil.GeneratedData("configuration/update", true), test.configurationPath)
					os.RemoveAll(projectName)
					if !test.projectDirectoryMissing {
						if err := os.MkdirAll(projectName, 0755); err != nil {
							t.Fatalf("\nCannot create test project: %v", err)
						}
					}
					if test.projectDirectoryIsFile {
						if err := os.MkdirAll(filepath.Dir(projectName), 0755); err != nil {
							t.Fatalf("\nCannot create test project: %v", err)
						}
						file, _ := os.OpenFile(projectName, os.O_RDONLY|os.O_CREATE, 0666)
						file.Close()
					}
				}
			}

			err := Update(configurationPath, test.environment, projectName, test.serviceName)

			testutil.VerifyError(test.expectedError, err, t)

			expected := testutil.ReadFile(path.Join(testData, test.expectedConfigurationFile), t, err)
			actual := testutil.ReadFile(path.Join(projectName, "configurations", test.serviceName, "values.yaml"), t, err)
			if err == nil && actual != expected {
				t.Fatalf("\nUnexpected configuration file:\nExpected: %v\nGot: %v", expected, actual)
			}
		})
	}
}
