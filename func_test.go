package main

import (
	`crypto/sha256`
	`fmt`
	`io/ioutil`
	`path/filepath`
	`reflect`
	`strings`
	`testing`
	`github.com/jscherff/gotest`
)

/* Test each function of the application by calling the function directly
   and comparing expected results with actual results.

	Config Helper Functions:

	[ ] newConfig(cf string) (*Config, error)
	[ ] loadConfig(t interface{}, cf string) error
	[ ] makePath(path string) (string, error)

	Router Functions:

	[X] route(i interface{}) (err error)
	[X] convert(i interface{}) (interface{}, error)
	[ ] update(i interface{}) (interface{})

	Action Functions:

	[ ] audit(dev usb.Auditer) (err error)
	[X] report(dev usb.Reporter) (err error)
	[X] serial(dev usb.Serializer) (err error)

	API Client Functions:

	[X] newSn(dev usb.Serializer) (string, error)
	[X] checkin(dev usb.Reporter) (error)
	[X] checkout(dev usb.Auditer) ([]byte, error)
	[ ] sendAudit(dev usb.Auditer) (error)
	[ ] vendor(dev usb.Updater) (s string, err error)
	[ ] product(dev usb.Updater) (s string, err error)

	HTTP Helper Functions:

	[ ] httpPost(url string, j []byte ) (b []byte, hs httpStatus, err error)
	[ ] httpGet(url string) (b []byte, hs httpStatus, err error)
	[ ] httpRequest(req *http.Request) (b []byte, hs httpStatus, err error)
*/

func TestFuncSerial(t *testing.T) {

	var err error

	t.Run("serial() and newSn() Must Obtain Serial Number (Magtek)", func(t *testing.T) {

		resetFlags(t)
		td.Mag[`mag1`].SerialNum = ``

		td.Mag[`mag1`].SerialNum, err = newSn(td.Mag[`mag1`])
		gotest.Ok(t, err)
		gotest.Assert(t, td.Mag[`mag1`].SerialNum != ``, `empty SN provided by server`)
		//TODO: assert correct serial number format
	})

	t.Run("serial() and newSn() Must Obtain Serial Number (IDTech)", func(t *testing.T) {

		resetFlags(t)
		td.Idt[`idt1`].SerialNum = ``

		td.Idt[`idt1`].SerialNum, err = newSn(td.Idt[`idt1`])
		gotest.Ok(t, err)
		gotest.Assert(t, td.Idt[`idt1`].SerialNum != ``, `empty SN provided by server`)
		//TODO: assert correct serial number format
	})

	t.Run("serial() and newSn() Must Obtain Serial Number (Unknown)", func(t *testing.T) {

		resetFlags(t)
		td.Idt[`idt1`].SerialNum = ``
		td.Idt[`idt1`].ObjectType = `*usb.Unknown`

		td.Idt[`idt1`].SerialNum, err = newSn(td.Idt[`idt1`])
		gotest.Ok(t, err)
		gotest.Assert(t, td.Idt[`idt1`].SerialNum != ``, `empty SN provided by server`)
		//TODO: assert correct serial number format
	})


	restoreState(t)
}

func TestFuncReport(t *testing.T) {

	var err error

	t.Run("(*Device).CSV() Must Match SHA256 Signature", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `csv`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := fmt.Sprintf(`%s-V%s-P%s.%s`,
			td.Mag[`mag1`].Conn(),
			td.Mag[`mag1`].VID(),
			td.Mag[`mag1`].PID(),
			*fReportFormat,
		)
		fn = filepath.Join(conf.Paths.ReportDir, fn)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`CSV`][`mag1`], `unexpected hash signature of CSV report`)
	})

	t.Run("(*Device).NVP() Must Match SHA256 Signature", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `nvp`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := fmt.Sprintf(`%s-V%s-P%s.%s`,
			td.Mag[`mag1`].Conn(),
			td.Mag[`mag1`].VID(),
			td.Mag[`mag1`].PID(),
			*fReportFormat,
		)
		fn = filepath.Join(conf.Paths.ReportDir, fn)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`NVP`][`mag1`], `unexpected hash signature of NVP report`)
	})

	t.Run("(*Device).XML() Must Match SHA256 Signature", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `xml`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := fmt.Sprintf(`%s-V%s-P%s.%s`,
			td.Mag[`mag1`].Conn(),
			td.Mag[`mag1`].VID(),
			td.Mag[`mag1`].PID(),
			*fReportFormat,
		)
		fn = filepath.Join(conf.Paths.ReportDir, fn)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`PXML`][`mag1`], `unexpected hash signature of XML report`)
	})

	t.Run("(*Device).JSON() Must Match SHA256 Signature", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `json`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := fmt.Sprintf(`%s-V%s-P%s.%s`,
			td.Mag[`mag1`].Conn(),
			td.Mag[`mag1`].VID(),
			td.Mag[`mag1`].PID(),
			*fReportFormat,
		)
		fn = filepath.Join(conf.Paths.ReportDir, fn)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`PJSN`][`mag1`], `unexpected hash signature of JSON report`)
	})
}

// Test device checkin and checkout.
func TestFuncCheckInOut(t *testing.T) {

	var err error

	t.Run("checkin() and checkout() Devices Must Match", func(t *testing.T) {

		resetFlags(t)

		err = checkin(td.Mag[`mag1`])
		gotest.Ok(t, err)

		j, err := checkout(td.Mag[`mag1`])
		gotest.Ok(t, err)

		ss, err := td.Mag[`mag1`].CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) == 0, `unmodified device should match last checkin`)
	})

	t.Run("checkin() and checkout() Devices Must Not Match", func(t *testing.T) {

		resetFlags(t)

		err = checkin(td.Mag[`mag1`])
		gotest.Ok(t, err)

		td.Mag[`mag1`].SoftwareID = `21042818B02`

		j, err := checkout(td.Mag[`mag1`])
		gotest.Ok(t, err)

		ss, err := td.Mag[`mag1`].CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) != 0, `modified device should not match last checkin`)
	})

	restoreState(t)
}

func TestFuncAudit(t *testing.T) {

	var err error

	t.Run("audit() Must Show Changes", func(t *testing.T) {

		resetFlags(t)
		*fActionAudit = true

		err = checkin(td.Mag[`mag1`])
		gotest.Ok(t, err)

		err = audit(td.Mag[`mag1`])
		gotest.Ok(t, err)

		gotest.Assert(t, len(td.Mag[`mag1`].Changes) == 0, `device change log should be empty`)

		err = audit(td.Mag[`mag2`])
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(td.Mag[`mag2`].Changes, td.Chg),
			`device change log does not contain known device differences`)

		fb, err := ioutil.ReadFile(conf.Loggers.Logger[`change`].LogFile)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, td.Clg[0]) && strings.Contains(fs, td.Clg[1]),
			`application change log does not contain known device differences`)
	})

	restoreState(t)
}
