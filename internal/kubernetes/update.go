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
	"github.com/callsign/gitops/internal/directory"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// Update the Kubernetes resources
func Update(projectName, serviceName, environment string) error {

	if projectName == "" {
		return fmt.Errorf("Missing project name")
	}
	if serviceName == "" {
		return fmt.Errorf("Missing service name")
	}
	if environment == "" {
		return fmt.Errorf("Missing environment")
	}

	temporaryDirectory, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("Cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(temporaryDirectory)

	if err := executeHelmTemplate(projectName, environment, temporaryDirectory); err != nil {
		return fmt.Errorf("Cannot execute helm template: %v", err)
	}

	if err := copyResources(projectName, serviceName, temporaryDirectory); err != nil {
		return fmt.Errorf("Cannot copy Kubernetes resources: %v", err)
	}

	return nil
}

func executeHelmTemplate(projectName, environment, temporaryDirectory string) error {

	var packagedChartDirectoryEntries []os.FileInfo
	var err error
	if packagedChartDirectoryEntries, err = ioutil.ReadDir("build/packages/helm"); err != nil {
		return fmt.Errorf("Cannot read build/packages/helm directory")
	}
	if len(packagedChartDirectoryEntries) == 0 {
		return fmt.Errorf("Missing packaged chart in build/packages/helm")
	}
	packagedChartFilename := packagedChartDirectoryEntries[len(packagedChartDirectoryEntries) - 1].Name()

	chartArgument := fmt.Sprintf("build/packages/helm/%s", packagedChartFilename)
	namespaceArgument := fmt.Sprintf("--namespace=%s-%s", projectName, environment)
	outputDirectoryArgument := fmt.Sprintf("--output-dir=%s", temporaryDirectory)
	valuesArgument := fmt.Sprintf("environments/%s/values.yaml", environment)

	command := exec.Command("helm", "template", chartArgument, namespaceArgument, outputDirectoryArgument, "-f", valuesArgument)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", output)
	}
	return nil
}

func copyResources(projectName, serviceName, temporaryDirectory string) error {

	projectPath, _ := filepath.Abs(projectName)
	var projectInfo os.FileInfo
	var err error
	if projectInfo, err = os.Stat(projectPath); err != nil {
		return fmt.Errorf("Project directory does not exist (%s)", projectPath)
	}
	if !projectInfo.IsDir() {
		return fmt.Errorf("Project directory is a file (%s)", projectPath)
	}

	source, _ := filepath.Abs(path.Join(temporaryDirectory, serviceName, "templates"))
	destination := path.Join(projectPath, "kubernetes", serviceName)

	return directory.Copy(source, destination)
}
