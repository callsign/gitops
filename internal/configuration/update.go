// Copyright © 2018 Callsign. All rights reserved.

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
