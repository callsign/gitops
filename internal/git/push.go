// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"fmt"
)

// Push the GitOps project changes
func Push(projectName, commitMessage string) error {
	fmt.Println("Pushing GitOps project changes")
	if err := git(projectName, "add", "."); err != nil {
		return err
	}
	if err := git(projectName, "commit", "-m", commitMessage); err != nil {
		return err
	}
	if err := git(projectName, "push"); err != nil {
		return err
	}
	return nil
}
