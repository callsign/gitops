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
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

// RequirementsYaml a Helm chart requirements.yaml
type RequirementsYaml struct {
	Dependencies []Dependency
}

// Dependency a Helm chart dependency
type Dependency struct {
	Name       string
	Version    string
	Repository string
	Alias      string   `yaml:",omitempty"`
	Condition  string   `yaml:",omitempty"`
	Tags       []string `yaml:",omitempty"`
}

const stableHelmRepository = "https://kubernetes-charts.storage.googleapis.com"
var companyHelmRepository = "https://artifactory.a2org.com/artifactory/helm-callsign-release"

// Update the service in the GitOps project by
// 1) updating the version in the Helm requirements.yaml file
// 2) updating the Helm requirements.lock file
func Update(projectName, serviceName, serviceVersion string) error {

	if serviceName == "" {
		return fmt.Errorf("Missing service name")
	}
	if serviceVersion == "" {
		return fmt.Errorf("Missing service version")
	}

	var chartPath string
	var err error
	if chartPath, err = getChartPath(projectName); err != nil {
		return fmt.Errorf("Cannot find chart: %v", err)
	}

	var helmRepositories []string
	if helmRepositories, err = updateRequirementsYaml(chartPath, serviceName, serviceVersion); err != nil {
		return fmt.Errorf("Cannot update requirements.yaml: %v", err)
	}

	if err = addHelmRepositories(helmRepositories); err != nil {
		return fmt.Errorf("Cannot add Helm repositories: %v", err)
	}

	if err = updateRequirementsLock(chartPath); err != nil {
		return fmt.Errorf("Cannot update requirements.lock: %v", err)
	}

	return nil
}

func getChartPath(projectName string) (string, error) {
	chartParent := projectName + "/charts"
	if _, err := os.Stat(chartParent); err != nil {
		return "", fmt.Errorf("Missing charts directory")
	}

	var chartParentEntries []os.FileInfo
	var err error
	if chartParentEntries, err = ioutil.ReadDir(chartParent); err != nil {
		return "", fmt.Errorf("Cannot read charts directory")
	}

	if len(chartParentEntries) == 0 || !chartParentEntries[0].IsDir() {
		return "", fmt.Errorf("Missing chart")
	}

	return path.Join(chartParent, chartParentEntries[0].Name()), nil
}

func updateRequirementsYaml(chartPath string, serviceName, serviceVersion string) ([]string ,error) {
	requirementsYaml := RequirementsYaml{}
	requirementsYamlPath := path.Join(chartPath, "/requirements.yaml")
	if _, err := os.Stat(requirementsYamlPath); err == nil {
		requirementsYamlFile, err := ioutil.ReadFile(requirementsYamlPath)
		if err != nil {
			return nil, fmt.Errorf("Cannot read file: %v", err)
		}
		if err := yaml.Unmarshal(requirementsYamlFile, &requirementsYaml); err != nil {
			return nil, fmt.Errorf("Cannot unmarshal file: %v", err)
		}
	}

	var serviceFound bool
	helmRepositories := make([]string, 0)
	for index, dependency := range requirementsYaml.Dependencies {
		if dependency.Name == serviceName {
			dependency.Version = serviceVersion
			requirementsYaml.Dependencies[index] = dependency
			serviceFound = true
		}
		helmRepositories = addHelmRepository(dependency.Repository, helmRepositories)
	}

	if !serviceFound {
		dependency := Dependency{Name: serviceName, Version: serviceVersion, Repository: companyHelmRepository}
		requirementsYaml.Dependencies = append(requirementsYaml.Dependencies, dependency)
		helmRepositories = addHelmRepository(companyHelmRepository, helmRepositories)
	}

	updatedFile, err := yaml.Marshal(&requirementsYaml)
	if err != nil {
		return nil, fmt.Errorf("Cannot marshal updated file")
	}
	if err := ioutil.WriteFile(path.Join(chartPath, "requirements.yaml"), updatedFile, 0644); err != nil {
		return nil, fmt.Errorf("Cannot write updated file")
	}

	return helmRepositories, nil
}

func addHelmRepository(repository string, repositories []string) []string {
	if repository != stableHelmRepository &&
	   !stringInSlice(repository, repositories) &&
	   !strings.HasPrefix(repository, "file") {
		return append(repositories, repository)
	}
	return repositories
}

func addHelmRepositories(repositories []string) error {
	for _, repository := range repositories {
		random := make([]byte, 4)
	    rand.Read(random)
		name := hex.EncodeToString(random)
		command := exec.Command("helm", "repo", "add", name, repository)
		if output, err := command.CombinedOutput(); err != nil {
			return fmt.Errorf("%s", output)
		}
	}
	return nil
}

func updateRequirementsLock(chartPath string) error {
	command := exec.Command("helm", "dependency", "update")
	command.Dir = chartPath
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", output)
	}
	return nil
}

func stringInSlice(element string, list []string) bool {
    for _, current := range list {
        if current == element {
            return true
        }
    }
    return false
}
