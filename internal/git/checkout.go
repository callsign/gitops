// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Checkout the GitOps project
func Checkout(projectURL, environment string) (string, error) {

	fmt.Println("Cloning GitOps project")
	if _, err := Git(".", "clone", projectURL); err != nil {
		return "", err
		//return "", fmt.Errorf("Cannot git clone %s: %s", projectURL, output)
	}

	projectName := removeExtension(filepath.Base(projectURL))

	fmt.Printf("Checking out %s environment branch\n", environment)
	if _, err := Git(projectName, "checkout", environment); err != nil {
		return "", err
	}

	return projectName, nil
}

func removeExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}
