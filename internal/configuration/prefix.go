// Copyright © 2018 Callsign. All rights reserved.

package configuration

import (
	"path/filepath"
    "fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

func prefix(projectName, serviceName string) error {
	configurationPath, _ :=  filepath.Abs(path.Join(projectName, "configurations", serviceName))

	if _, err := os.Stat(configurationPath); err != nil {
		return fmt.Errorf("Cannot find configuration (%s)", configurationPath)
	}

	var environmentEntries []os.FileInfo
	var err error
	if environmentEntries, err = ioutil.ReadDir(configurationPath); err != nil {
		return err
	}
	for _, environmentEntry := range environmentEntries {
		environmentPath := path.Join(configurationPath, environmentEntry.Name())

		var files []string
		filepath.Walk(environmentPath, func(path string, file os.FileInfo, _ error) error {
			if !file.IsDir() {
				if filepath.Ext(path) == ".yaml" {
					files = append(files, file.Name())
				}
			}
			return nil
		})

		for _, file := range files {
			filePath := path.Join(environmentPath, file)

			var fileBytes []byte
			var err error
			if fileBytes, err = ioutil.ReadFile(filePath); err != nil {
				return fmt.Errorf("Cannot read configuration file (%s)", filePath)
			}

			fileContent := make(map[interface{}]interface{})
			if err = yaml.Unmarshal([]byte(fileBytes), &fileContent); err != nil {
				return fmt.Errorf("Invalid YAML (%s)", filePath)
			}

			if len(fileContent) == 0 {
				continue;
			}

			updatedEntryContent := make(map[interface{}]interface{})
			updatedEntryContent[serviceName] = fileContent

			if fileBytes, err = yaml.Marshal(&updatedEntryContent); err != nil {
				return fmt.Errorf("Cannot marshal updated configuration file (%s)", filePath)
			}

			if err = ioutil.WriteFile(filePath, fileBytes, 0644); err != nil {
				return fmt.Errorf("Cannot write updated configuration file (%s)", filePath)
			}
		}
	}
	return nil
}
