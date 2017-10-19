// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	`encoding/json`
	`fmt`
	`path/filepath`
	`os`
)

var (
	program = filepath.Base(os.Args[0])
	version = `undefined`
)

// Config holds the application configuration settings. The struct tags
// must match the field names in the JSON configuration file.
type Config struct {

	Paths struct {
		AppDir          string
		LogDir          string
		StateDir        string
		ReportDir       string
	}

	Files struct {
		SystemLog       string
		ChangeLog       string
		ErrorLog        string
		Legacy          string
	}

	Server struct {
		URL             string
		Endpoint        map[string]string
	}


	Logging struct {

		System struct {
			Logfile bool
			Console bool
			Syslog  bool
		}

		Change struct {
			Logfile bool
			Console bool
			Syslog  bool
		}

		Error struct {
			Logfile bool
			Console bool
			Syslog  bool
		}
	}

	Syslog struct {
		Protocol        string
		Port            string
		Host            string
	}

	Include struct {
		VendorID        map[string]bool
		ProductID       map[string]map[string]bool
		Default         bool
	}

	Format struct {
		Report          string
		Default         string
	}

	DebugLevel              int
}

// newConfig retrieves the settings in the JSON configuration file and
// populates the fields in the runtime configuration. It also creates
// directories if they do not already exist.
func newConfig(cf string) (this *Config, err error) {

	this = &Config{}
	ad := filepath.Dir(os.Args[0])

	if dn := filepath.Dir(cf); len(dn) == 0 {
		cf = filepath.Join(ad, cf)
	}

	fh, err := os.Open(cf)

	if err != nil {
		return nil, err
	}

	defer fh.Close()
	jd := json.NewDecoder(fh)

	if err = jd.Decode(&this); err != nil {
		return nil, err
	}

	this.Paths.AppDir = ad

	// Helpers to prepend and/or create paths as necessary.

	var mkd = func(pd, d string) (string, error) {
		if dn := filepath.Dir(d); dn == `` || dn == `.` {
			d = filepath.Join(pd, d)
		}
		return d, os.MkdirAll(d, DirMode)
	}

	var mkf = func(pd, f string) (string, error) {

		if dn := filepath.Dir(f); dn == `` || dn == `.` {
			f = filepath.Join(pd, f)
			return f, os.MkdirAll(pd, DirMode)
		} else {
			return f, os.MkdirAll(dn, DirMode)
		}
	}

	// Build directory names and create paths as necessary. If a directory
	// is relative, prepend the application directory.

	if this.Paths.LogDir, err = mkd(this.Paths.AppDir, this.Paths.LogDir); err != nil {
		return nil, err
	}
	if this.Paths.StateDir, err = mkd(this.Paths.AppDir, this.Paths.StateDir); err != nil {
		return nil, err
	}
	if this.Paths.ReportDir, err = mkd(this.Paths.AppDir, this.Paths.ReportDir); err != nil {
		return nil, err
	}

	// Build file names and create paths as necessary. If a filename is 
	// relative, prepend the appropriate application directory.

	if this.Files.SystemLog, err = mkf(this.Paths.LogDir, this.Files.SystemLog); err != nil {
		return nil, err
	}
	if this.Files.ChangeLog, err = mkf(this.Paths.LogDir, this.Files.ChangeLog); err != nil {
		return nil, err
	}
	if this.Files.ErrorLog, err = mkf(this.Paths.LogDir, this.Files.ErrorLog); err != nil {
		return nil, err
	}
	if this.Files.Legacy, err = mkf(this.Paths.AppDir, this.Files.Legacy); err != nil {
		return nil, err
	}

	return this, nil
}

// displayVersion displays the program version.
func displayVersion() {
        fmt.Fprintf(os.Stderr, "%s version %s\n", program, version)
}
