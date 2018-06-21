// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Git executes a GIT command
func Git(directory string, arguments  ...string) (string, error) {
	command := exec.Command("git", arguments...)
	command.Dir = directory
	outputbytes, err := command.CombinedOutput()
	output := strings.Trim(string(outputbytes[:]), "\n")
	if err != nil {
		return "", fmt.Errorf("Cannot git %s: %v", arguments[0], output)
	}
	return output, nil
}
