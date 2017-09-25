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
	`log`
	`os`
	`path/filepath`
	`reflect`
	`strings`
	`testing`
	`github.com/jscherff/gocmdb/usbci`
	`github.com/jscherff/gotest`
)

func init() {

	magChanges[0] = []string{`SoftwareID`, `21042818B01`, `21042818B02`}
	magChanges[1] = []string{`USBSpec`, `1.10`, `2.00`}

	var err, errM1, errM2, errG1, errG2 error

	if mag1, errM1 = usbci.NewMagtek(nil); errM1 == nil {
		errM1 = mag1.RestoreJSON(mag1JSON)
	}

	if mag2, errM2 = usbci.NewMagtek(nil); errM2 == nil {
		errM2 = mag2.RestoreJSON(mag2JSON)
	}

	if gen1, errG1 = usbci.NewGeneric(nil); errG1 == nil {
		errG1 = gen1.RestoreJSON(gen1JSON)
	}

	if gen2, errG2 = usbci.NewGeneric(nil); errG2 == nil {
		errG2 = gen2.RestoreJSON(gen2JSON)
	}

	if errM1 != nil || errM2 != nil || errG1 != nil || errG2 != nil {
		log.Fatal(os.Stderr, "Testing setup failed: could not restore devices.")
	}

	if conf, err = newConfig(`config.json`); err != nil {
		os.Exit(1)
		log.Fatal(os.Stderr, "Testing setup failed: could not restore configuration.")
	}

	conf.Logging.Console = false
	slog, clog, elog = newLoggers()
}

func TestGetNewSN(t *testing.T) {

	t.Run("GetNewSN() Function", func(t *testing.T) {

		j, err := mag2.JSON()
		gotest.Ok(t, err)

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreJSON(j)
		gotest.Ok(t, err)

		mag3.SerialNum = ``

		s, err := getNewSN(mag3)
		gotest.Ok(t, err)
		gotest.Assert(t, len(s) != 0, `empty SN provided by server`)

		s, err = getNewSN(mag2)
		gotest.Assert(t, err != nil, `request for SN when device has one should produce error`)
	})

	t.Run("serialAction() Function", func(t *testing.T) {
		//TODO
	})
}

func TestReporting(t *testing.T) {

	*fReportConsole = false

	t.Run("JSON Report", func(t *testing.T) {

		*fReportFormat = `json`
		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err := reportAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1SigPrettyJSON == sha256.Sum256(b), `unexpected hash signature of JSON report`)
	})

	t.Run("XML Report", func(t *testing.T) {

		*fReportFormat = `xml`
		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err := reportAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1SigPrettyXML == sha256.Sum256(b), `unexpected hash signature of XML report`)
	})

	t.Run("CSV Report", func(t *testing.T) {

		*fReportFormat = `csv`
		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err := reportAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1SigCSV == sha256.Sum256(b), `unexpected hash signature of CSV report`)
	})

	t.Run("NVP Report", func(t *testing.T) {

		*fReportFormat = `nvp`
		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err := reportAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1SigNVP == sha256.Sum256(b), `unexpected hash signature of NVP report`)
	})

	t.Run("Legacy Report", func(t *testing.T) {

		*fActionLegacy = true

		err := legacyAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(conf.Files.Legacy)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1SigLegacy == sha256.Sum256(b), `unexpected hash signature of Legacy report`)
	})

}

// Test device checkin and checkout.
func TestCheckinCheckout(t *testing.T) {

	t.Run("Checkin and Checkout Must Match", func(t *testing.T) {

		err := checkinDevice(mag1)
		gotest.Ok(t, err)

		j, err := checkoutDevice(mag1)
		gotest.Ok(t, err)

		ss, err := mag1.CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) == 0, `unmodified device should match last checkin`)
	})

	t.Run("Checkin and Checkout Must Not Match", func(t *testing.T) {

		mag1mod, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag1mod.RestoreJSON(mag1JSON)
		gotest.Ok(t, err)

		mag1mod.SoftwareID = `21042818B02`
		err = checkinDevice(mag1mod)
		gotest.Ok(t, err)

		j, err := checkoutDevice(mag1mod)
		gotest.Ok(t, err)

		ss, err := mag1.CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) != 0, `modified device should not match last checkin`)
	})
}

func TestAuditing(t *testing.T) {

	var (
		ch1 = `"SoftwareID" was "21042818B01", now "21042818B02"`
		ch2 = `"USBSpec" was "1.10", now "2.00"`
	)


	t.Run("Local Audit", func(t *testing.T) {

		*fAuditLocal = true
		*fAuditServer = false

		err := auditAction(mag1)
		gotest.Assert(t, err != nil, `first run should result in file-not-found error`)

		err = auditAction(mag1)
		gotest.Ok(t, err)

		gotest.Assert(t, len(mag1.Changes) == 0, `device change log should be empty`)

		mag1mod, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag1mod.RestoreJSON(mag2JSON)
		gotest.Ok(t, err)

		err = auditAction(mag1mod)
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(mag1mod.Changes, magChanges),
			`device change log does not contain known device differences`)

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, ch1) && strings.Contains(fs, ch2),
			`application change log does not contain known device differences`)
	})

	t.Run("Server Audit", func(t *testing.T) {

		*fAuditLocal = false
		*fAuditServer = true

		err := checkinDevice(mag1)
		gotest.Ok(t, err)

		err = auditAction(mag1)
		gotest.Ok(t, err)

		gotest.Assert(t, len(mag1.Changes) == 0, `device change log should be empty`)

		mag1mod, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag1mod.RestoreJSON(mag2JSON)
		gotest.Ok(t, err)

		err = auditAction(mag1mod)
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(mag1mod.Changes, magChanges),
			`device change log does not contain known device differences`)

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, ch1) && strings.Contains(fs, ch2),
			`application change log does not contain known device differences`)
	})
}

func TestFileReadWrite(t *testing.T) {

	var (
		b []byte
		err error
	)

	// File Write Paths

	wfn1 := `test1.txt`
	wfn2 := `log/test2.txt`
	wfn3 := filepath.Join(os.Getenv(`TEMP`), `test3.txt`)

	// File Read Paths ('should')

	rfn1 := filepath.Join(conf.Paths.AppDir, `test1.txt`)
	rfn2 := `log/test2.txt`
	rfn3 := filepath.Join(os.Getenv(`TEMP`), `test3.txt`)

	// Generate file content

	j, err := mag1.JSON()
	gotest.Ok(t, err)

	// File Write Tests

	err = writeFile(j, wfn1)
	gotest.Ok(t, err)

	err = writeFile(j, wfn2)
	gotest.Ok(t, err)

	err = writeFile(j, wfn3)
	gotest.Ok(t, err)

	// File Read Tests

	b, err = readFile(rfn1)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == mag1SigJSON, `unexpected hash signature of file contents`)

	b, err = readFile(rfn2)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == mag1SigJSON, `unexpected hash signature of file contents`)

	b, err = readFile(rfn3)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == mag1SigJSON, `unexpected hash signature of file contents`)

	// File Read Test Validations

	b, err = ioutil.ReadFile(rfn1)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == mag1SigJSON, `unexpected hash signature of file contents`)

	b, err = ioutil.ReadFile(rfn2)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == mag1SigJSON, `unexpected hash signature of file contents`)

	b, err = ioutil.ReadFile(rfn3)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == mag1SigJSON, `unexpected hash signature of file contents`)
}



/*

TODO:
	serialAction(o gocmdb.Configurable) (err error)
	httpPost(url string, j []byte ) (b []byte, sc int, err error)
	httpGet(url string) (b []byte, sc int, err error)
	httpRequest(req *http.Request) (b []byte, sc int, err error)
	newLoggers() (sl, cl, el *log.Logger)
	magtekRouter(musb gocmdb.MagtekUSB) (err error)
	genericRouter(gusb gocmdb.GenericUSB) (err error)

WIP:


DONE:
	newConfig(string) (*Config, error) - init()
	getNewSN(o gocmdb.Registerable) (string, error) - TestGetNewSN()
	reportAction(o gocmdb.Reportable) (error) - TestReporting()
	legacyAction(o gocmdb.Reportable) (error) - TestReporting()
	checkinDevice(o gocmdb.Registerable) (error) - TestCheckinCheckout()
	checkoutDevice(o gocmdb.Auditable) ([]byte, error) - TestCheckinCheckout()
	auditAction(o gocmdb.Auditable) (error) - TestAuditing()
	submitAudit(o gocmdb.Auditable) (error) - TestAuditing()
	readFile(string, []byte) (error) - TestFileReadWrite()
	writeFile([]byte, string) (error) - TestFileReadWrite()

*/
