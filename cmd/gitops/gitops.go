// Copyright © 2018 Callsign. All rights reserved.

package main

import (
	"bytes"
	"log"
	"os"

	"callsign.com/gitops/internal/configuration"
	"callsign.com/gitops/internal/deployment"
	"callsign.com/gitops/internal/service"
)

func main() {
	arguments := os.Args[1:]
	log.SetFlags(0)
	if len(arguments) == 2 && arguments[0] == "request-deployment" {
		if err := deployment.Request(arguments[1]); err != nil {
			log.Fatalln(err)
		}
	} else if len(arguments) == 3 && arguments[0] == "update-configuration" {
		if err := configuration.Update("environments", arguments[1], arguments[2], service.Name()); err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println(usage())
	}
}

func usage() string {
	var buffer bytes.Buffer
	buffer.WriteString("Usage: gitops <command> <argument(s)>\n")
	buffer.WriteString("       Valid commands:\n")
	buffer.WriteString("          gitops request-deployment <project-url>\n")
	buffer.WriteString("          gitops update-configuration <environment> <project-path>")
	return buffer.String()
}
