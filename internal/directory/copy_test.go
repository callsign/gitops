// Copyright © 2018 Callsign. All rights reserved.

package directory

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/callsign/gitops/internal/testutil"
)

func Test_Copy(t *testing.T) {
	destination := testutil.GeneratedData("directory/copy", true)
	var err error
	if err = Copy(testutil.Data("directory/copy"), destination); err != nil {
		t.Fatalf("\nUnexpected error: %v", err)
	}
	if _, err = os.Stat(destination); err != nil {
		t.Fatalf("\nMissing destination directory")
	}
	var destinationEntries []os.FileInfo
	if destinationEntries, err = ioutil.ReadDir(destination); err != nil {
		t.Fatalf("\nError reading destination directory: %v", err)
	}
	if len(destinationEntries) != 2 {
		t.Fatalf("\nWrong number of destination directory entries")
	}

	recursivelyCopiedDiretory := path.Join(destination, "directory")
	if _, err = os.Stat(recursivelyCopiedDiretory); err != nil {
		t.Fatalf("\nMissing recursively copied directory")
	}
	var recursivelyCopiedDiretoryEntries []os.FileInfo
	if recursivelyCopiedDiretoryEntries, err = ioutil.ReadDir(recursivelyCopiedDiretory); err != nil {
		t.Fatalf("\nError reading recursively copied directory: %v", err)
	}
	if len(recursivelyCopiedDiretoryEntries) != 1 {
		t.Fatalf("\nWrong number of recursively copied directory entries")
	}
}
