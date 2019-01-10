// Copyright © 2018 Callsign. All rights reserved.

package environment

import (
	"os"
	"fmt"
	"regexp"

	"github.com/callsign/gitops/internal/git"
)

var environments = map[string]string{
	"production": "^master$",
	"staging":    "^(release|hotfix)\\/\\S+$",
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

	for environment, serviceBranchRegexp := range environments {
		if match, _ := regexp.MatchString(serviceBranchRegexp, serviceBranch); match {
			return environment, nil
		}
	}

	return "", fmt.Errorf("Cannot find deployment environment for service branch %s", serviceBranch)
}
