/*
 * Copyright 2012-2018 the original author or authors.
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

// CustomDeploymentsYaml a custom deployments Yaml
type CustomDeploymentsYaml struct {
	Deployments []CustomDeployment
}

// CustomDeployment a custom deployment
type CustomDeployment struct {
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
	return serviceBranchToEnvironment(serviceBranch, "custom-deployments.yaml")
}

func serviceBranchToEnvironment(serviceBranch, customDeploymentsYamlPath string) (string, error) {

	for environment, serviceBranchRegexp := range environments {
		if match, _ := regexp.MatchString(serviceBranchRegexp, serviceBranch); match {
			return environment, nil
		}
	}

	customDeploymentsYaml := CustomDeploymentsYaml{}
	if _, err := os.Stat(customDeploymentsYamlPath); err == nil {
		customDeploymentsYamlFile, err := ioutil.ReadFile(customDeploymentsYamlPath)
		if err != nil {
			return "", fmt.Errorf("Cannot read %s: %v", customDeploymentsYamlPath, err)
		}
		if err := yaml.Unmarshal(customDeploymentsYamlFile, &customDeploymentsYaml); err != nil {
			return "", fmt.Errorf("Cannot unmarshal %s: %v", customDeploymentsYamlPath, err)
		}
		for _, customDeployment := range customDeploymentsYaml.Deployments {
			if customDeployment.Branch == serviceBranch {
				for standardEnvironment := range environments {
					if customDeployment.Environment == standardEnvironment {
						format := "Deployment of service branch %s to deployment environment %s forbidden"
						return "", fmt.Errorf(format, serviceBranch, standardEnvironment)
					}
				}
				return customDeployment.Environment, nil
			}
		}
	}

	return "", fmt.Errorf("Cannot find deployment environment for service branch %s", serviceBranch)
}
