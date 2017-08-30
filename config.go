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
	"encoding/json"
	"path/filepath"
	"fmt"
	"os"
)

// The filename of the JSON configuration file.
const configFile string = "config.json"

// The config variable holds the runtime configuration.
var config *Config

// init retrieves the settings in the JSON configuration file and
// populates the fields in the runtime configuration.
func init() {

	ap := fmt.Sprintf("%s%c", filepath.Dir(os.Args[0]), filepath.Separator)
	fh, e := os.Open(ap + configFile)

	defer fh.Close()

	if e == nil {
		jd := json.NewDecoder(fh)
		e = jd.Decode(&c)
	}

	if e == nil {
		c.AppPath = ap
	}

	if e != nil {
		log.Fatalf("error reading config: %v", e)
	}
}

// Config holds the application configuration settings. The struct tags
// must match the field names in the JSON configuration file.
type Config struct {
	AppDir		string				`json:"app_dir"`

	LogDir		string				`json:"log_dir"`
	OutDir		string				`json:"out_dir"`
	AuditDir	string				`json:"audit_dir"`

	LogFile		string				`json:"log_file"`
	OutFile		string				`json:"out_file"`
	LegacyLogFile	string				`json:"legacy_log_file"`
	LegacyOutFile	string				`json:"legacy_out_file"`

	AuditUrl	string				`json:"audit_url"`
	CheckinUrl	string				`json:"checkin_url"`

	IncludeVID	map[string]bool			`json:"include_vid"`
	IncludePID	map[string]map[string]bool	`json:"include_pid"`

	DefaultInclude	bool				`json:"default_include"`
	DefaultFormat	string				`json:"default_format"`
}
