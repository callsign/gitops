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

package environment

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/callsign/gitops/internal/git"
	"gopkg.in/yaml.v2"
)

// DeploymentsYaml a deployments Yaml
type DeploymentsYaml struct {
	Deployments []eployment
}

// Deployment a deployment
type Deployment struct {
	Branch      string
	Environment string
}

var environments = map[string]string{
	"prod":    "^master$",
	"staging": "^(release|hotfix)\\/\\S+$",
	"dev":     "^develop$",
}

// Get the environment to deploy to
func Get() (string, error) {
	serviceBranch, err := git.Git(".", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("Cannot read service branch: %s", serviceBranch)
	}
	if serviceBranch == "HEAD" {
		serviceBranch = os.Getenv("CI_COMMIT_REF_NAME")
	}
	return serviceBranchToEnvironment(serviceBranch, "deployments.yaml")
}

func serviceBranchToEnvironment(serviceBranch, deploymentsYamlPath string) (string, error) {

	for environment, serviceBranchRegexp := range environments {
		if match, _ := regexp.MatchString(serviceBranchRegexp, serviceBranch); match {
			return environment, nil
		}
	}

	deploymentsYaml := DeploymentsYaml{}
	if _, err := os.Stat(deploymentsYamlPath); err == nil {
		deploymentsYamlFile, err := ioutil.ReadFile(deploymentsYamlPath)
		if err != nil {
			return "", fmt.Errorf("Cannot read %s: %v", deploymentsYamlPath, err)
		}
		if err := yaml.Unmarshal(deploymentsYamlFile, &deploymentsYaml); err != nil {
			return "", fmt.Errorf("Cannot unmarshal %s: %v", deploymentsYamlPath, err)
		}
		for _, deployment := range deploymentsYaml.Deployments {
			if deployment.Branch == serviceBranch {
				for standardEnvironment := range environments {
					if deployment.Environment == standardEnvironment {
						format := "Deployment of service branch %s to deployment environment %s forbidden"
						return "", fmt.Errorf(format, serviceBranch, standardEnvironment)
					}
				}
				return deployment.Environment, nil
			}
		}
	}

	return "", fmt.Errorf("Cannot find deployment environment for service branch %s", serviceBranch)
}
