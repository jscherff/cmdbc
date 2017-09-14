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
	`io/ioutil`
	`os`
	`strings`
	`github.com/RackSec/srslog`
)

const (
	PriInfo = srslog.LOG_LOCAL7|srslog.LOG_INFO
	PriErr = srslog.LOG_LOCAL7|srslog.LOG_ERR
	FileFlags = os.O_APPEND|os.O_CREATE|os.O_WRONLY
	FileMode = 0640
	DirMode = 0750
)

func NewLoggers() (sl, cl, el *log.Logger) {

	var sw, cw, ew []io.Writer

	var newf func(f string) (h *os.File, err error) {

		if h, err = os.OpenFile(f, FileFlags, FileMode); err != nil {
			log.Printf(`%v`, ErrorDecorator(err))
		}

		return h, err
	}

	var newsl func(prot, raddr, tag string, pri int) (s *srslog.Writer, err error) {

		if sl, err = srslog.Dial(prot, raddr, pri, tag); err != nil {
			log.Printf(`%v`, ErrorDecorator(err))
		}
		return s, err
	}

	if conf.Logging.LogFiles {

		if f, err := newf(conf.Files.SystemLog); err == nil {
			sw = append(sw, f)
		}
		if f, err := newf(conf.Files.ChangeLog); err == nil {
			cw = append(cw, f)
		}
		if f, err := newf(conf.Files.ErrorLog); err == nil {
			ew = append(ew, f)
		}
	}

	if conf.Logging.Console {
		sw = append(sw, os.Stdout)
		cw = append(cw, os.Stdout)
		ew = append(ew, os.Stderr)
	}

	if conf.Logging.Syslog {

		port, prot, addr := conf.Syslog.Port, conf.Syslog.Protocol, conf.Syslog.Address
		raddr := strings.Join([]string{addr, port}, `:`)

		if s, err := newsl(port, prot, addr, "gocmdbcli", PriInfo); err == nil {
			sw = append(sw, s)
		}
		if s, err := newsl(port, prot, addr, "gocmdbcli", PriInfo); err == nil {
			cw = append(cw, s)
		}
		if s, err := newsl(port, prot, addr, "gocmdbcli", PriErr); err == nil {
			ew = append(ew, s)
		}
	}

	if len(sw) == 0 {
		sw = append(sw, ioutil.Discard)
	}
	if len(cw) == 0 {
		cw = append(cw, ioutil.Discard)
	}
	if len(ew) == 0 {
		ew = append(ew, ioutil.Discard)
	}

	sl = log.New(io.MultiWriter(sw...), `system: `, log.LstdFlags)
	cl = log.New(io.MultiWriter(cw...), `change: `, log.LstdFlags)
	el = log.New(io.MultiWriter(ew...), `error: `,  log.LstdFlags)

	return this, err
}
