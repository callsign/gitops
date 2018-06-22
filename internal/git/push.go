// Copyright © 2018 Callsign. All rights reserved.

package git

// Push the GitOps project changes
func Push(projectName, commitMessage string) error {
	if _, err := Git(".", "config", "--global", "user.name", "GitLab CI"); err != nil {
		return err
	}
	if _, err := Git(".", "config", "--global", "user.email", "gitlab-ci@callsign.com"); err != nil {
		return err
	}
	if _, err := Git(projectName, "add", "."); err != nil {
		return err
	}
	if _, err := Git(projectName, "commit", "-m", commitMessage); err != nil {
		return err
	}
	if _, err := Git(projectName, "push"); err != nil {
		return err
	}
	return nil
}
