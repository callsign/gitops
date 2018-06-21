// Copyright © 2018 Callsign. All rights reserved.

package deployment

import (
	"fmt"
	"testing"

	"github.com/callsign/gitops/internal/testutil"
)

func Test_Push(t *testing.T) {
	tests := []struct {
		name          string
		projectURL    string
		expectedError error
	}{{
		name:          "should return an error on invalid project URL",
		projectURL:    "invalid",
		expectedError: fmt.Errorf("Cannot git clone: *"),
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualError := Request(test.projectURL)
			testutil.VerifyError(test.expectedError, actualError, t)
		})
	}
}
