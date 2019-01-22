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

package configuration

import (
	"fmt"
)

// Update the service configuration in the GitOps project
// 1) Copy environments/* to configurations/<service>/
// 2) Prefix the configuration keys with the service name
func Update(configurationPath, projectName, serviceName string) error {

	if err := copy(configurationPath, projectName, serviceName); err != nil {
		return fmt.Errorf("Cannot copy service configuration: %v", err)
	}

	if err := prefix(projectName, serviceName); err != nil {
		return fmt.Errorf("Cannot add service prefix to service configuration: %v", err)
	}

	return nil
}
