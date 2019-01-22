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
