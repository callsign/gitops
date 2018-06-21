// Copyright © 2018 Callsign. All rights reserved.

package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name        string
		arguments   []string
		expectedLog string
	}{{
		name: "should log usage when no command",
		arguments: []string{"gitops"},
	}, {
		name: "should log usage when invalid command",
		arguments: []string{"gitops", "invalid-command"},
	}, {
		name: "should log usage when missing request-deployment parameter",
		arguments: []string{"gitops", "request-deployment"},
	}, {
		name: "should log usage when too many request-deployment parameters",
		arguments: []string{"gitops", "request-deployment", "https://project-url", "invalid-parameter"},
	}, {
		name: "should log usage when missing update-configuration parameters",
		arguments: []string{"gitops", "update-configuration"},
	}, {
		name: "should log usage when too many update command parameters",
		arguments: []string{"gitops", "update-configuration", "staging", "../project", "invalid-parameter"},
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.arguments
			var logs bytes.Buffer
			log.SetOutput(&logs)
			main()
			actualLogs := strings.TrimSuffix(logs.String(), "\n")
			expectedLogs := usage()
			if (actualLogs != expectedLogs) {
				t.Fatalf("\nUnexpected result:\nExpected: %v\nGot: %v", expectedLogs, actualLogs)
			}
		})
	}
}
