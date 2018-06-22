// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"path/filepath"
	"strings"
)

// Checkout the GitOps project
func Checkout(projectURL, environment string) (string, error) {
	if _, err := Git(".", "clone", projectURL); err != nil {
		return "", err
	}
	projectName := removeExtension(filepath.Base(projectURL))
	if _, err := Git(projectName, "checkout", environment); err != nil {
		return "", err
	}
	return projectName, nil
}

func removeExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}
