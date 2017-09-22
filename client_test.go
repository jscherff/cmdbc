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
	`crypto/sha256`
	`io/ioutil`
	`os`
	`path/filepath`
	`strings`
	`testing`
	`github.com/jscherff/gocmdb/usbci`
	`github.com/jscherff/gotest`
)

func init() {

	var err error

	if conf, err = NewConfig(`config.json`); err != nil {
		os.Exit(1)
	}

	conf.Logging.Console = false
	slog, clog, elog = NewLoggers()
}

// Test reporting functions.
func TestReporting(t *testing.T) {

	var (
		b []byte
		fn string
		err error
	)

	*fReportConsole = false

	t.Run("JSN Report", func(t *testing.T) {

		*fReportFormat = `json`
		fn = filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err = reportAction(mag1)
		gotest.Ok(t, err)

		b, err = ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1ShaJSN == sha256.Sum256(b), `unexpected hash signature of JSON report`)
	})

	t.Run("XML Report", func(t *testing.T) {

		*fReportFormat = `xml`
		fn = filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err = reportAction(mag1)
		gotest.Ok(t, err)

		b, err = ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1ShaXML == sha256.Sum256(b), `unexpected hash signature of XML report`)
	})

	t.Run("CSV Report", func(t *testing.T) {

		*fReportFormat = `csv`
		fn = filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err = reportAction(mag1)
		gotest.Ok(t, err)

		b, err = ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1ShaCSV == sha256.Sum256(b), `unexpected hash signature of CSV report`)
	})

	t.Run("NVP Report", func(t *testing.T) {

		*fReportFormat = `nvp`
		fn = filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err = reportAction(mag1)
		gotest.Ok(t, err)

		b, err = ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1ShaNVP == sha256.Sum256(b), `unexpected hash signature of NVP report`)
	})
}

// Test device checkin and checkout.
func TestCheckinCheckout(t *testing.T) {

	var (
		j []byte
		err error
		ss [][]string
		mag1mod *usbci.Magtek
	)

	t.Run("Checkin Checkout Must Match", func(t *testing.T) {

		err = CheckinDevice(mag1)
		gotest.Ok(t, err)

		j, err = CheckoutDevice(mag1)
		gotest.Ok(t, err)

		ss, err = mag1.CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) == 0, `unmodified device does not match self`)
	})

	t.Run("Checkin Checkout Must Not Match", func(t *testing.T) {

		mag1mod, err = usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag1mod.RestoreJSON(j)
		gotest.Ok(t, err)

		mag1mod.SoftwareID = `21042818B02`
		err = CheckinDevice(mag1mod)
		gotest.Ok(t, err)

		j, err = CheckoutDevice(mag1mod)
		gotest.Ok(t, err)

		ss, err = mag1.CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) != 0, `modified device matches self`)
	})
}

func TestGetNewSN(t *testing.T) {

	var (
		s string
		err error
	)

	j, err := mag2.JSON()
	gotest.Ok(t, err)

	mag3, err := usbci.NewMagtek(nil)
	gotest.Ok(t, err)

	err = mag3.RestoreJSON(j)
	gotest.Ok(t, err)

	mag3.SerialNum = ``

	s, err = GetNewSN(mag3)
	gotest.Ok(t, err)
	gotest.Assert(t, len(s) != 0, `empty SN provided by server`)

	s, err = GetNewSN(mag2)
	gotest.Assert(t, err != nil, `request for SN when device has one should produce error`)
}

func TestAuditAction(t *testing.T) {

	var (
		err error
		ch1 = `"SoftwareID" was "21042818B01", now "21042818B02"`
		ch2 = `"USBSpec" was "1.10", now "2.00"`
	)


	t.Run("Local Audit", func(t *testing.T) {

		*fAuditLocal = true

		err = auditAction(mag2)
		gotest.Assert(t, err != nil, `first run should result in file-not-found error`)

		err = auditAction(mag1)
		gotest.Ok(t, err)

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)

		gotest.Assert(t, strings.Contains(fs, ch1) && strings.Contains(fs, ch2),
			`known device differences not recorded in change log`)
	})
}
