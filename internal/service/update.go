// Copyright © 2018 Callsign. All rights reserved.

package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

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
	Alias      string
	Condition  string
	Tags       []string
}

// Update the service in the GitOps project by
// 1) updating the version in the Helm requirements.yaml file
// 2) updating the Helm requirements.lock file
func Update(projectName, serviceName, serviceVersion string) error {

	fmt.Println("Updating requirements.yaml")
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

	var requirementsYaml []byte
	if requirementsYaml, err = readRequirementsYaml(chartPath); err != nil {
		return fmt.Errorf("Cannot read requirements.yaml: %v", err)
	}

	if err = updateRequirementsYaml(requirementsYaml, serviceName, serviceVersion, chartPath); err != nil {
		return fmt.Errorf("Cannot update requirements.yaml: %v", err)
	}

	fmt.Println("Updating requirements.lock")
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

	if len(chartParentEntries) == 0 {
		return "", fmt.Errorf("Missing chart")
	}

	if len(chartParentEntries) > 1 {
		return "", fmt.Errorf("Multiple charts")
	}

	return path.Join(chartParent, chartParentEntries[0].Name()), nil
}

func readRequirementsYaml(chartPath string) ([]byte, error) {
	path := path.Join(chartPath, "/requirements.yaml")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("Missing requirements.yaml")
	}
	return ioutil.ReadFile(path)
}

func updateRequirementsYaml(file []byte, serviceName, serviceVersion, chartPath string) error {
	requirementsYaml := RequirementsYaml{}
	if err := yaml.Unmarshal(file, &requirementsYaml); err != nil {
		return fmt.Errorf("Cannot unmarshal file")
	}

	var serviceFound bool
	for index, dependency := range requirementsYaml.Dependencies {
		if dependency.Name == serviceName {
			dependency.Version = serviceVersion
			requirementsYaml.Dependencies[index] = dependency
			serviceFound = true
		}
	}

	if !serviceFound {
		return fmt.Errorf("Missing service entry")
	}

	updatedFile, err := yaml.Marshal(&requirementsYaml)
	if err != nil {
		return fmt.Errorf("Cannot marshal updated file")
	}

	if err := ioutil.WriteFile(path.Join(chartPath, "requirements.yaml"), updatedFile, 0644); err != nil {
		return fmt.Errorf("Cannot write updated file")
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
