// Copyright © 2018 Callsign. All rights reserved.

package environment

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var environments = map[string]string{
	"production": "^master$",
	"staging":    "^(release|hotfix)\\/\\S+$",
	"verify":     "^develop$",
}

// Get the environment to deploy to
func Get() (string, error) {
	fmt.Println("Determiniing deployment environment")

	output, err := exec.Command("git", "symbolic-ref", "--short", "HEAD").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Cannot read service branch: %s", output)
	}
	serviceBranch := strings.Trim(string(output[:]), "\n")

	for environment, serviceBranchRegexp := range environments {
		if match, _ := regexp.MatchString(serviceBranchRegexp, serviceBranch); match {
			return environment, nil
		}
	}

	return "", fmt.Errorf("Cannot find deployment environment for service branch %s", serviceBranch)
}
