// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"path/filepath"
	"io/ioutil"
	"os"
)

func writeFile(b []byte, p string) (e error) {

	d, f := filepath.Split(p)

	if len(d) == 0 {
		d = config.AppPath
		p = filepath.Join(d, f)
	}

	if e = os.MkdirAll(d, 0755); e == nil {
		e = ioutil.WriteFile(p, b, 0644)
	}

	return e
}

func readFile(p string, b []byte) (e error) {

	d, f := filepath.Split(p)

	if len(d) == 0 {
		d = config.AppPath
		p = filepath.Join(d, f)
	}

	b, e = ioutil.ReadFile(p)

	return e
}
