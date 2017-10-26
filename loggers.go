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
	`log`
	`io`
	`io/ioutil`
	`os`
	`path/filepath`
)

const (
        LogFileAppend = os.O_APPEND|os.O_CREATE|os.O_WRONLY
        LogFileMode = 0640
)


var (
	LogFlags = map[string]int{
		`date`:		log.Ldate,
		`time`:		log.Ltime,
		`file`:		log.Lshortfile,
	}
)

// Loggers contains a collection of log.Logger objects with embedded
// configuraiton and multiwriter capabilities.
type Loggers struct {
	LogDir string
	Logger map[string]*Logger
	Console bool
	Syslog bool
}

// Init initializes each Logger with embedded properties and parameters.
func (this *Loggers) Init(syslog *Syslog) error {

	if dn, err := makePath(this.LogDir); err != nil {
		return err
	} else {
		this.LogDir = dn
	}

	for tag, logger := range this.Logger {

		tag += ` `

		logger.LogFile = filepath.Join(this.LogDir, logger.LogFile)
		logger.Console = logger.Console || this.Console
		logger.Syslog = logger.Syslog || this.Syslog

		if err := logger.Init(tag, syslog); err != nil {
			return err
		}
	}

	return nil
}

// Logger is a log.Logger object with embedded configuration and
// multiwriter capabilities.
type Logger struct {
	*log.Logger
	LogFile string
	Console bool
	Syslog bool
	Prefix []string
}

// Init initializes the Logger with embedded properties and parameters.
func (this *Logger) Init(tag string, syslog *Syslog) error {

	var (
		writers []io.Writer
		flags int
	)

	if file, err := os.OpenFile(this.LogFile, LogFileAppend, LogFileMode); err != nil {
		return err
	} else {
		writers = append(writers, file)
	}

	if this.Console {
		writers = append(writers, os.Stdout)
	}

	if this.Syslog && syslog != nil {
		writers = append(writers, syslog)
	}

	if len(writers) == 0 {
		writers = append(writers, ioutil.Discard)
	}

	for _, flag := range this.Prefix {
		flags |= LogFlags[flag]
	}

	this.Logger = log.New(io.MultiWriter(writers...), tag, flags)

	return nil
}
