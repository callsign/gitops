// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"fmt"
	"os/exec"
)

func git(directory string, arguments  ...string) error {
	command := exec.Command("git", arguments...)
	command.Dir = directory
	if err := command.Run(); err != nil {
		return fmt.Errorf("Cannot git %s: %v", arguments[0], err)
	}
	return nil
}
