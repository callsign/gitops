// Copyright © 2018 Callsign. All rights reserved.

package service

import (
	"path/filepath"	
)

// Name returns the service name
func Name() (string) {
	path, _ := filepath.Abs(".")
	return filepath.Base(path)
}
