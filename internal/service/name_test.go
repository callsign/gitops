// Copyright © 2018 Callsign. All rights reserved.

package service

import (
	"testing"
)

func Test_Name(t *testing.T) {
	tests := []struct {
		name string
		expected string
	}{{		
		name: "should return the service name",
		expected: "service",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Name()
			if (actual != test.expected) {
				t.Fatalf("\nUnexpected result:\nExpected: %v\nGot: %v", test.expected, actual)
			}
		})
	}
}
