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
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Cannot git %s: %v", arguments[0], output)
	}
	return strings.Trim(string(output[:]), "\n"), nil
}
