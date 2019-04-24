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
	    filePath := path.Join(configurationPath, environmentEntry.Name(), "values.yaml")
        if _, err := os.Stat(filePath); err != nil {
            continue
        }

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
	return nil
}
