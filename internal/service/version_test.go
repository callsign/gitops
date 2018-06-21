// Copyright © 2018 Callsign. All rights reserved.

package service

import (
	"fmt"
	"path"
	"testing"

	"github.com/callsign/gitops/internal/testutil"
)

func Test_version(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedError  error
		expectedResult string
	}{{
		name:          "should return a error if the version file is missing",
		input:         "missing-file",
		expectedError: fmt.Errorf("Cannot read service version file (*"),
	}, {
		name:           "should return the version",
		input:          "valid-version",
		expectedResult: "2.0.1805200912+abcde",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			servicePath := path.Join(testutil.Data("service/version"), test.input)

			actualResult, actualError := version(servicePath)

			testutil.VerifyError(test.expectedError, actualError, t)

			if actualResult != test.expectedResult {
				t.Fatalf("\nUnexpected result:\nExpected: %v\nGot: |%v|", test.expectedResult, actualResult)
			}
		})
	}
}

func Test_Version(t *testing.T) {
	tests := []struct {
		name          string
		expectedError error
	}{{
		name:          "should return a error if the version file is missing",
		expectedError: fmt.Errorf("Cannot read service version file (*"),
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, actualError := Version()
			testutil.VerifyError(test.expectedError, actualError, t)
		})
	}
}
