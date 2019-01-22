/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
