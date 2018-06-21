// Copyright © 2018 Callsign. All rights reserved.

package environment

import (
	"testing"
)

func Test_Get(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{{
		name: "should return the environment",
		expected: "production",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Get()
			if (err != nil) {
				t.Fatalf("\nTest failure: %v", err)
			}
			if (actual != test.expected) {
				t.Fatalf("\nUnexpected result:\nExpected: %v\nGot: %v", test.expected, actual)
			}
		})
	}
}
