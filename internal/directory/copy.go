// Copyright © 2018 Callsign. All rights reserved.

package directory

import (
	"fmt"
    "io/ioutil"
    "os"
    "path"
)

// Copy copied a directory recursively
func Copy(source string, destination string) error {
	var err error
	var sourceInfo os.FileInfo
	var sourceEntries []os.FileInfo

	if sourceInfo, err = os.Stat(source); err != nil {
		return fmt.Errorf("Source directory does not exist (%s)", source)
	}

	if !sourceInfo.IsDir() {
		return fmt.Errorf("Source directory is a file (%s)", source)
	}
	
	if err = os.MkdirAll(destination, sourceInfo.Mode()); err != nil {
		return err
	}

	if sourceEntries, err = ioutil.ReadDir(source); err != nil {
		return err
	}
	for _, sourceEntry := range sourceEntries {
		sourceEntryPath := path.Join(source, sourceEntry.Name())
		destinationEntryPath := path.Join(destination, sourceEntry.Name())

		if sourceEntry.IsDir() {
			if err = Copy(sourceEntryPath, destinationEntryPath); err != nil {
				return err
			}
		} else {
			if err = copyFile(sourceEntryPath, destinationEntryPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(source string, destination string) error {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(source); err != nil {
		return err
	}
	if err = ioutil.WriteFile(destination, data, 0644); err != nil {
		return err
	}
	return nil
}
