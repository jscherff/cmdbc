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
	`fmt`
	`io/ioutil`
	`os`
	`path/filepath`
	`reflect`
	`strings`
	`testing`
	`github.com/jscherff/gotest`
)

func TestFuncSerial(t *testing.T) {

	var err error

	t.Run("GetNewSN() Function", func(t *testing.T) {

		resetFlags(t)
		mag1.SerialNum = ``

		mag1.SerialNum, err = getNewSN(mag1)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1.SerialNum != ``, `empty SN provided by server`)
	})

	restoreState(t)
}

func TestFuncReport(t *testing.T) {

	var err error

	t.Run("JSON Report", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `json`

		err = reportAction(mag1)
		gotest.Ok(t, err)

		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == mag1SigPJSON, `unexpected hash signature of JSON report`)
	})

	t.Run("XML Report", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `xml`

		err = reportAction(mag1)
		gotest.Ok(t, err)

		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == mag1SigPXML, `unexpected hash signature of XML report`)
	})

	t.Run("CSV Report", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `csv`

		err = reportAction(mag1)
		gotest.Ok(t, err)

		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == mag1SigCSV, `unexpected hash signature of CSV report`)
	})

	t.Run("NVP Report", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `nvp`

		err = reportAction(mag1)
		gotest.Ok(t, err)

		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == mag1SigNVP, `unexpected hash signature of NVP report`)
	})

	t.Run("Legacy Report", func(t *testing.T) {

		resetFlags(t)
		*fActionLegacy = true

		err = legacyAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(conf.Files.Legacy)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == mag1SigLegacy, `unexpected hash signature of Legacy report`)
	})
}

// Test device checkin and checkout.
func TestFuncCheckInOut(t *testing.T) {

	var err error

	t.Run("Checkin and Checkout Must Match", func(t *testing.T) {

		resetFlags(t)

		err = checkinDevice(mag1)
		gotest.Ok(t, err)

		j, err := checkoutDevice(mag1)
		gotest.Ok(t, err)

		ss, err := mag1.CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) == 0, `unmodified device should match last checkin`)
	})

	t.Run("Checkin and Checkout Must Not Match", func(t *testing.T) {

		resetFlags(t)

		err = checkinDevice(mag1)
		gotest.Ok(t, err)

		mag1.SoftwareID = `21042818B02`

		j, err := checkoutDevice(mag1)
		gotest.Ok(t, err)

		ss, err := mag1.CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) != 0, `modified device should not match last checkin`)
	})

	restoreState(t)
}

func TestFuncAudit(t *testing.T) {

	var err error

	t.Run("Local Audit", func(t *testing.T) {

		resetFlags(t)
		*fAuditLocal = true

		af := fmt.Sprintf(`%s-%s-%s.json`, mag1.VID(), mag1.PID(), mag1.ID())
		os.RemoveAll(filepath.Join(conf.Paths.StateDir, af))

		err = auditAction(mag1)
		gotest.Assert(t, err != nil, `first run should result in file-not-found error`)

		err = auditAction(mag1)
		gotest.Ok(t, err)

		gotest.Assert(t, len(mag1.Changes) == 0, `device change log should be empty`)

		err = auditAction(mag2)
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(mag2.Changes, magChanges),
			`device change log does not contain known device differences`)

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, ClogCh1) && strings.Contains(fs, ClogCh2),
			`application change log does not contain known device differences`)
	})

	t.Run("Server Audit", func(t *testing.T) {

		resetFlags(t)
		*fAuditServer = true

		err = checkinDevice(mag1)
		gotest.Ok(t, err)

		err = auditAction(mag1)
		gotest.Ok(t, err)

		gotest.Assert(t, len(mag1.Changes) == 0, `device change log should be empty`)

		err = auditAction(mag2)
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(mag2.Changes, magChanges),
			`device change log does not contain known device differences`)

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, ClogCh1) && strings.Contains(fs, ClogCh2),
			`application change log does not contain known device differences`)
	})

	restoreState(t)
}

func TestFuncFileIO(t *testing.T) {

	// File Write Paths

	wfn1 := `test1.txt`
	wfn2 := `log/test2.txt`
	wfn3 := filepath.Join(os.Getenv(`TEMP`), `test3.txt`)

	// File Read Paths ('should')

	rfn1 := filepath.Join(conf.Paths.AppDir, `test1.txt`)
	rfn2 := `log/test2.txt`
	rfn3 := filepath.Join(os.Getenv(`TEMP`), `test3.txt`)

	// Generate file content

	b, err := mag1.JSON()
	gotest.Ok(t, err)

	// File Write Tests

	err = writeFile(b, wfn1)
	gotest.Ok(t, err)

	err = writeFile(b, wfn2)
	gotest.Ok(t, err)

	err = writeFile(b, wfn3)
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

	os.RemoveAll(wfn1)
	os.RemoveAll(wfn2)
	os.RemoveAll(wfn3)
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
