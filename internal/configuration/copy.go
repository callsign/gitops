// Copyright © 2018 Callsign. All rights reserved.

package configuration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/callsign/gitops/internal/directory"
)

func copy(configurationPath, environment, projectName, serviceName string) error {
	if configurationPath == "" {
		return fmt.Errorf("Configuration path missing")
	}
	if environment == "" {
		return fmt.Errorf("Environment missing")
	}
	if projectName == "" {
		return fmt.Errorf("Project name missing")
	}
	if serviceName == "" {
		return fmt.Errorf("Service name missing")
	}

	projectPath, _ := filepath.Abs(projectName)
	var projectInfo os.FileInfo
	var err error
	if projectInfo, err = os.Stat(projectPath); err != nil {
		return fmt.Errorf("Project directory does not exist (%s)", projectPath)
	}
	if !projectInfo.IsDir() {
		return fmt.Errorf("Project directory is a file (%s)", projectPath)
	}

	source, _ := filepath.Abs(path.Join(configurationPath, environment))
	destination := path.Join(projectPath, "configurations", serviceName)

	return directory.Copy(source, destination)
}
