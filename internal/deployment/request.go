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

package deployment

import (
	"fmt"
	"path/filepath"

	"github.com/callsign/gitops/internal/configuration"
	"github.com/callsign/gitops/internal/environment"
	"github.com/callsign/gitops/internal/git"
	"github.com/callsign/gitops/internal/kubernetes"
	"github.com/callsign/gitops/internal/service"
)

// Request a service deployment
func Request(projectURL string) error {
	environment, err := environment.Get()
	if err != nil {
		return err
	}

	projectName, err := git.Checkout(projectURL, environment)
	if err != nil {
		return err
	}

	serviceName := service.Name()
	err = kubernetes.Update(projectName, serviceName, environment)
	if err != nil {
		return err
	}	
	
	configurationPath, _ := filepath.Abs("environments")
	err = configuration.Update(configurationPath, projectName, serviceName)
	if err != nil {
		return err
	}

	serviceVersion, err := service.Version()
	if err != nil {
		return err
	}

	fmt.Printf("Requesting deployment of %s %s to %s\n", serviceName, serviceVersion, environment)

	err = service.Update(projectName, serviceName, serviceVersion)
	if err != nil {
		return err
	}

	commitMessage := serviceName + "-" + serviceVersion
	err = git.Push(projectName, commitMessage)
	if err != nil {
		return err
	}

	return nil
}
