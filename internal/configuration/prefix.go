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

	var files []string
	filepath.Walk(configurationPath, func(path string, file os.FileInfo, _ error) error {
		if !file.IsDir() {
			if filepath.Ext(path) == ".yaml" {
				files = append(files, file.Name())
			}
		}
		return nil
	})

	for _, file := range files {
		filePath := path.Join(configurationPath, file)

		var fileBytes []byte
		var err error
		if fileBytes, err = ioutil.ReadFile(filePath); err != nil {
			return fmt.Errorf("Cannot read configuration file (%s)", filePath)
		}

		fileContent := make(map[interface{}]interface{})
		if err = yaml.Unmarshal([]byte(fileBytes), &fileContent); err != nil {
			return fmt.Errorf("Invalid YAML (%s)", filePath)
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
	return nil
}
