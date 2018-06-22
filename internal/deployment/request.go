// Copyright © 2018 Callsign. All rights reserved.

package deployment

import (
	"fmt"
	"path/filepath"

	"github.com/callsign/gitops/internal/configuration"
	"github.com/callsign/gitops/internal/environment"
	"github.com/callsign/gitops/internal/git"
	"github.com/callsign/gitops/internal/service"
)

// Request a service deployment
func Request(projectURL string) error {
	environment, err := environment.Get()
	if err != nil {
		return err
	}

	projectName, err := git.Checkout(projectURL, environment)
	if err != nil {
		return err
	}

	serviceName := service.Name()
	configurationPath, _ := filepath.Abs("environments")
	err = configuration.Update(configurationPath, environment, projectName, serviceName)
	if err != nil {
		return err
	}

	serviceVersion, err := service.Version()
	if err != nil {
		return err
	}

	fmt.Printf("Requesting deployment of %s %s to %s\n", serviceName, serviceVersion, environment)

	err = service.Update(projectName, serviceName, serviceVersion)
	if err != nil {
		return err
	}

	commitMessage := serviceName + "-" + serviceVersion
	err = git.Push(projectName, commitMessage)
	if err != nil {
		return err
	}

	return nil
}
