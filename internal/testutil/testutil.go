// Copyright © 2018 Callsign. All rights reserved.

package testutil

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// Data returns the test data
func Data(test string) string {
	path, _ := filepath.Abs(path.Join("../../test-data", test))
	return path
}

// GeneratedData returns the test generated data
func GeneratedData(tested string, absolute bool) string {
	path := filepath.Join("../../test-data/generated", tested)
	if absolute {
		path, _ = filepath.Abs(path)
	}
	return path
}

// ReadFile reads a file if their is no previous error
func ReadFile(path string, t *testing.T, previousError error) string {
	var err error
	var bytes []byte
	if previousError != nil {
		return ""
	}
	if bytes, err = ioutil.ReadFile(path); err != nil {
		t.Fatalf("\nCannot read file %v: %v", path, err)
	}
	return string(bytes)
}

// VerifyError verifies that the actual error is the expected error
func VerifyError(expectedError, actualError error, t *testing.T) {
	if expectedError != nil && actualError == nil {
		t.Fatalf("\nUnexpected error:\nExpected: %v\nGot: %v", expectedError, actualError)
	}
	if expectedError == nil && actualError != nil {
		t.Fatalf("\nUnexpected error:\nExpected: %v\nGot: %v", expectedError, actualError)
	}
	if expectedError != nil && actualError != nil {
		if strings.HasSuffix(expectedError.Error(), "*") {
			if !strings.HasPrefix(actualError.Error(), strings.TrimSuffix(expectedError.Error(), "*")) {
				t.Fatalf("\nUnexpected error:\nExpected: %v\nGot: %v", expectedError, actualError)
			}
		} else if expectedError.Error() != actualError.Error() {
			t.Fatalf("\nUnexpected error:\nExpected: %v\nGot: %v", expectedError, actualError)
		}
	}
}
