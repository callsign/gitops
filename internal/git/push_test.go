// Copyright © 2018 Callsign. All rights reserved.

package git

import (
	"fmt"
	"testing"

	"github.com/callsign/gitops/internal/testutil"
)

func Test_Push(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		expectedError error
	}{{
		name:          "should return an error on invalid project name",
		projectName:   "invalid",
		expectedError: fmt.Errorf("Cannot git add: *"),
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualError := Push(test.projectName, "message")
			testutil.VerifyError(test.expectedError, actualError, t)
		})
	}
}
