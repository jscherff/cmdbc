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
		td.Mag[`mag1`].SerialNum = ``

		td.Mag[`mag1`].SerialNum, err = usbCiNewSnV1(td.Mag[`mag1`])
		gotest.Ok(t, err)
		gotest.Assert(t, td.Mag[`mag1`].SerialNum != ``, `empty SN provided by server`)
	})

	restoreState(t)
}

func TestFuncReport(t *testing.T) {

	var err error

	t.Run("JSON Report", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `json`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := filepath.Join(conf.Paths.ReportDir, td.Mag[`mag1`].Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`PJSN`][`mag1`], `unexpected hash signature of JSON report`)
	})

	t.Run("XML Report", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `xml`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := filepath.Join(conf.Paths.ReportDir, td.Mag[`mag1`].Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`PXML`][`mag1`], `unexpected hash signature of XML report`)
	})

	t.Run("CSV Report", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `csv`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := filepath.Join(conf.Paths.ReportDir, td.Mag[`mag1`].Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`CSV`][`mag1`], `unexpected hash signature of CSV report`)
	})

	t.Run("NVP Report", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `nvp`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := filepath.Join(conf.Paths.ReportDir, td.Mag[`mag1`].Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`NVP`][`mag1`], `unexpected hash signature of NVP report`)
	})
}

// Test device checkin and checkout.
func TestFuncCheckInOut(t *testing.T) {

	var err error

	t.Run("Checkin and Checkout Must Match", func(t *testing.T) {

		resetFlags(t)

		err = usbCiCheckinV1(td.Mag[`mag1`])
		gotest.Ok(t, err)

		j, err := usbCiCheckoutV1(td.Mag[`mag1`])
		gotest.Ok(t, err)

		ss, err := td.Mag[`mag1`].CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) == 0, `unmodified device should match last checkin`)
	})

	t.Run("Checkin and Checkout Must Not Match", func(t *testing.T) {

		resetFlags(t)

		err = usbCiCheckinV1(td.Mag[`mag1`])
		gotest.Ok(t, err)

		td.Mag[`mag1`].SoftwareID = `21042818B02`

		j, err := usbCiCheckoutV1(td.Mag[`mag1`])
		gotest.Ok(t, err)

		ss, err := td.Mag[`mag1`].CompareJSON(j)
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

		af := fmt.Sprintf(`%s-%s-%s.json`, td.Mag[`mag1`].VID(), td.Mag[`mag1`].PID(), td.Mag[`mag1`].ID())
		os.RemoveAll(filepath.Join(conf.Paths.StateDir, af))

		err = audit(td.Mag[`mag1`])
		gotest.Assert(t, err != nil, `first run should result in file-not-found error`)

		err = audit(td.Mag[`mag1`])
		gotest.Ok(t, err)

		gotest.Assert(t, len(td.Mag[`mag1`].Changes) == 0, `device change log should be empty`)

		err = audit(td.Mag[`mag2`])
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(td.Mag[`mag2`].Changes, td.Chg),
			`device change log does not contain known device differences`)

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, td.Clg[0]) && strings.Contains(fs, td.Clg[1]),
			`application change log does not contain known device differences`)
	})

	t.Run("Server Audit", func(t *testing.T) {

		resetFlags(t)
		*fAuditServer = true

		err = usbCiCheckinV1(td.Mag[`mag1`])
		gotest.Ok(t, err)

		err = audit(td.Mag[`mag1`])
		gotest.Ok(t, err)

		gotest.Assert(t, len(td.Mag[`mag1`].Changes) == 0, `device change log should be empty`)

		err = audit(td.Mag[`mag2`])
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(td.Mag[`mag2`].Changes, td.Chg),
			`device change log does not contain known device differences`)

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, td.Clg[0]) && strings.Contains(fs, td.Clg[1]),
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

	b, err := td.Mag[`mag1`].JSON()
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
	gotest.Assert(t, sha256.Sum256(b) == td.Sig[`JSN`][`mag1`], `unexpected hash signature of file contents`)

	b, err = readFile(rfn2)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == td.Sig[`JSN`][`mag1`], `unexpected hash signature of file contents`)

	b, err = readFile(rfn3)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == td.Sig[`JSN`][`mag1`], `unexpected hash signature of file contents`)

	// File Read Test Validations

	b, err = ioutil.ReadFile(rfn1)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == td.Sig[`JSN`][`mag1`], `unexpected hash signature of file contents`)

	b, err = ioutil.ReadFile(rfn2)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == td.Sig[`JSN`][`mag1`], `unexpected hash signature of file contents`)

	b, err = ioutil.ReadFile(rfn3)
	gotest.Ok(t, err)
	gotest.Assert(t, sha256.Sum256(b) == td.Sig[`JSN`][`mag1`], `unexpected hash signature of file contents`)

	os.RemoveAll(wfn1)
	os.RemoveAll(wfn2)
	os.RemoveAll(wfn3)
}

/*

TODO:
	serial(o gocmdb.Configurable) (err error)
	httpPost(url string, j []byte ) (b []byte, sc int, err error)
	httpGet(url string) (b []byte, sc int, err error)
	httpRequest(req *http.Request) (b []byte, sc int, err error)
	newLoggers() (sl, cl, el *log.Logger)
	magtekRouter(musb gocmdb.MagtekUSB) (err error)
	genericRouter(gusb gocmdb.GenericUSB) (err error)

DONE:
	newConfig(string) (*Config, error) - init()
	usbCiNewSnV1(o gocmdb.Registerable) (string, error) - TestGetNewSN()
	report(o gocmdb.Reportable) (error) - TestReporting()
	legacyHandler(o gocmdb.Reportable) (error) - TestReporting()
	usbCiCheckinV1(o gocmdb.Registerable) (error) - TestCheckinCheckout()
	usbCiCheckoutV1(o gocmdb.Auditable) ([]byte, error) - TestCheckinCheckout()
	audit(o gocmdb.Auditable) (error) - TestAuditing()
	usbCiAuditV1(o gocmdb.Auditable) (error) - TestAuditing()
	readFile(string, []byte) (error) - TestFileReadWrite()
	writeFile([]byte, string) (error) - TestFileReadWrite()

*/
