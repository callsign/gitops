// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"fmt"
	"os/exec"
)

// Push the GitOps project changes
func Push(projectName string) error {
	fmt.Println("Pushing GitOps project changes")
	command := exec.Command("git", "push")
	command.Dir = projectName
	if err := command.Run(); err != nil {
		return fmt.Errorf("Cannot push: %v", err)
	}
	return nil
}
