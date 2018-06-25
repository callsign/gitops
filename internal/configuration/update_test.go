// Copyright © 2018 Callsign. All rights reserved.

package configuration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/callsign/gitops/internal/testutil"
)

func Test_Update(t *testing.T) {
	tests := []struct {
		name                      string
		configurationPath         string
		projectNameMissing        bool
		projectDirectoryMissing   bool
		projectDirectoryIsFile    bool
		serviceName               string
		input                     string
		expectedError             error
		expectedConfigurationFile string
	}{{
		name:          "should return an error if the configuration directory is missing",
		serviceName:   "service",
		expectedError: fmt.Errorf("Cannot copy service configuration: Configuration path missing"),
	}, {
		name:              "should return an error if the configuration directory does not exist",
		configurationPath: "configuration-directory-does-not-exist",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Source directory does not exist (*"),
	}, {
		name:              "should return an error if the configuration directory is a file",
		configurationPath: "configuration-directory-is-file",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Source directory is a file (*"),
	}, {
		name:               "should return an error if the projectName is missing",
		configurationPath:  "project-name-missing",
		projectNameMissing: true,
		serviceName:        "service",
		expectedError:      fmt.Errorf("Cannot copy service configuration: Project name missing"),
	}, {
		name:                    "should return an error if the project directory does not exist",
		configurationPath:       "project-directory-does-not-exist",
		projectDirectoryMissing: true,
		serviceName:             "service",
		expectedError:           fmt.Errorf("Cannot copy service configuration: Project directory does not exist (*"),
	}, {
		name:                    "should return an error if the project directory is a file",
		configurationPath:       "project-directory-is-file",
		projectDirectoryMissing: true,
		projectDirectoryIsFile:  true,
		serviceName:             "service",
		expectedError:           fmt.Errorf("Cannot copy service configuration: Project directory is a file (*"),
	}, {
		name:              "should return an error if the serviceName is missing",
		configurationPath: "service-name-missing",
		expectedError:     fmt.Errorf("Cannot copy service configuration: Service name missing"),
	}, {
		name:              "should return an error on invalid yaml",
		configurationPath: "invalid-yaml",
		serviceName:       "service",
		expectedError:     fmt.Errorf("Cannot add service prefix to service configuration: Invalid YAML (*"),
	}, {
		name:                      "should copy the configuration and add a service prefix to the files",
		configurationPath:         "configuration-copied-service-prefix-added",
		serviceName:               "service",
		expectedConfigurationFile: "expected-configuration-file.yaml",
	}, {
		name:                      "should ignore non yaml files and directories",
		configurationPath:         "ignore-non-yaml-files-and-directories",
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

			err := Update(configurationPath, projectName, test.serviceName)

			testutil.VerifyError(test.expectedError, err, t)

			expected := testutil.ReadFile(path.Join(testData, test.expectedConfigurationFile), t, err)
			environments := []string{"staging", "production"}
			for _, environment := range environments {
				actual := testutil.ReadFile(path.Join(projectName, "configurations", test.serviceName, environment, "values.yaml"), t, err)
				if err == nil && actual != expected {
					t.Fatalf("\nUnexpected configuration file:\nExpected: %v\nGot: %v", expected, actual)
				}
			}
		})
	}
}
