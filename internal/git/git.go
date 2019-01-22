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
