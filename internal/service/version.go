// Copyright © 2018 Callsign. All rights reserved.

package service

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Version returns the service version
func Version() (string, error) {
	fmt.Println("Reading service version file")
	return version(".")
}

func version(servicePath string) (string, error) {
	var err error
	var bytes []byte
	path, _ := filepath.Abs(filepath.Join(servicePath, "build/packages/version"))
	if bytes, err = ioutil.ReadFile(path); err != nil {
		return "", fmt.Errorf("Cannot read service version file (%s): %v", path, err)
	}
	return strings.Trim(string(bytes), "\n"), nil
}
