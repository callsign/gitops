// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Checkout the GitOps project
func Checkout(projectURL, environment string) (string, error) {

	fmt.Println("Cloning GitOps project")
	if output, err := exec.Command("git", "clone", projectURL).CombinedOutput(); err != nil {
		return "", fmt.Errorf("Cannot git clone %s: %s", projectURL, output)
	}

	projectName := removeExtension(filepath.Base(projectURL))

	fmt.Printf("Checking out %s environment branch\n", environment)
	if err := git(projectName, "checkout", environment); err != nil {
		return "", err
	}

	return projectName, nil
}

func removeExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}
