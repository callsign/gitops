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

package service

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Version returns the service version
func Version() (string, error) {
	return version(".")
}

func version(servicePath string) (string, error) {
	var err error
	var bytes []byte
	path, _ := filepath.Abs(filepath.Join(servicePath, "build/packages/version"))
	if bytes, err = ioutil.ReadFile(path); err != nil {
		return "", fmt.Errorf("Cannot read service version file (%s): %v", path, err)
	}
	return strings.Trim(string(bytes), "\n"), nil
}
