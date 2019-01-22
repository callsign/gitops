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

package configuration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/callsign/gitops/internal/directory"
)

func copy(configurationPath, projectName, serviceName string) error {
	if configurationPath == "" {
		return fmt.Errorf("Configuration path missing")
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

	source, _ := filepath.Abs(configurationPath)
	destination := path.Join(projectPath, "configurations", serviceName)

	return directory.Copy(source, destination)
}
