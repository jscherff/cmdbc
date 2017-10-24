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

func TestFuncSerial(t *testing.T) {

	var err error

	t.Run("serial() and usbCiNewSNV1() Must Obtain Serial Number (Magtek)", func(t *testing.T) {

		resetFlags(t)
		td.Mag[`mag1`].SerialNum = ``

		td.Mag[`mag1`].SerialNum, err = usbCiNewSnV1(td.Mag[`mag1`])
		gotest.Ok(t, err)
		gotest.Assert(t, td.Mag[`mag1`].SerialNum != ``, `empty SN provided by server`)
		//TODO: assert correct serial number format
	})

	t.Run("serial() and usbCiNewSNV1() Must Obtain Serial Number (IDTech)", func(t *testing.T) {

		resetFlags(t)
		td.Idt[`idt1`].SerialNum = ``

		td.Idt[`idt1`].SerialNum, err = usbCiNewSnV1(td.Idt[`idt1`])
		gotest.Ok(t, err)
		gotest.Assert(t, td.Idt[`idt1`].SerialNum != ``, `empty SN provided by server`)
		//TODO: assert correct serial number format
	})

	t.Run("serial() and usbCiNewSNV1() Must Not Obtain Serial Number (Bad IDTech)", func(t *testing.T) {

		resetFlags(t)
		td.Idt[`idt1`].SerialNum = ``
		td.Idt[`idt1`].ObjectType = `*usb.Unknown`

		td.Idt[`idt1`].SerialNum, err = usbCiNewSnV1(td.Idt[`idt1`])
		gotest.Assert(t, err != nil, `attempt to obtain SN for unsupported device should fail`)
	})


	restoreState(t)
}

func TestFuncReport(t *testing.T) {

	var err error

	t.Run("(*Device).JSON() Must Match SHA256 Signature", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `json`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := fmt.Sprintf(`%s-%s.%s`, td.Mag[`mag1`].SN(), td.Mag[`mag1`].Conn(), *fReportFormat)
		fn = filepath.Join(conf.Paths.ReportDir, fn)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`PJSN`][`mag1`], `unexpected hash signature of JSON report`)
	})

	t.Run("(*Device).XML() Must Match SHA256 Signature", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `xml`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := fmt.Sprintf(`%s-%s.%s`, td.Mag[`mag1`].SN(), td.Mag[`mag1`].Conn(), *fReportFormat)
		fn = filepath.Join(conf.Paths.ReportDir, fn)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`PXML`][`mag1`], `unexpected hash signature of XML report`)
	})

	t.Run("(*Device).CSV() Must Match SHA256 Signature", func(t *testing.T) {

		resetFlags(t)
		*fReportFormat = `csv`

		err = report(td.Mag[`mag1`])
		gotest.Ok(t, err)

		fn := fmt.Sprintf(`%s-%s.%s`, td.Mag[`mag1`].SN(), td.Mag[`mag1`].Conn(), *fReportFormat)
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

		fn := fmt.Sprintf(`%s-%s.%s`, td.Mag[`mag1`].SN(), td.Mag[`mag1`].Conn(), *fReportFormat)
		fn = filepath.Join(conf.Paths.ReportDir, fn)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, sha256.Sum256(b) == td.Sig[`NVP`][`mag1`], `unexpected hash signature of NVP report`)
	})
}

// Test device checkin and checkout.
func TestFuncCheckInOut(t *testing.T) {

	var err error

	t.Run("usbCiCheckinV1() and usbCiCheckoutV1() Devices Must Match", func(t *testing.T) {

		resetFlags(t)

		err = usbCiCheckinV1(td.Mag[`mag1`])
		gotest.Ok(t, err)

		j, err := usbCiCheckoutV1(td.Mag[`mag1`])
		gotest.Ok(t, err)

		ss, err := td.Mag[`mag1`].CompareJSON(j)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss) == 0, `unmodified device should match last checkin`)
	})

	t.Run("usbCiCheckinV1() and usbCiCheckoutV1() Devices Must Not Match", func(t *testing.T) {

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

	t.Run("audit() Must Show Changes", func(t *testing.T) {

		resetFlags(t)
		*fActionAudit = true

		err = usbCiCheckinV1(td.Mag[`mag1`])
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

/*
TODO (OLD):
	[ ] serial(o gocmdb.Configurable) (err error)
	[ ] httpPost(url string, j []byte ) (b []byte, sc int, err error)
	[ ] httpGet(url string) (b []byte, sc int, err error)
	[ ] httpRequest(req *http.Request) (b []byte, sc int, err error)
	[ ] newLoggers() (sl, cl, el *log.Logger)
	[ ] magtekRouter(musb gocmdb.MagtekUSB) (err error)
	[ ] genericRouter(gusb gocmdb.GenericUSB) (err error)
	[X] newConfig(string) (*Config, error) - init()
	[X] usbCiNewSnV1(o gocmdb.Registerable) (string, error) - TestGetNewSN()
	[X] report(o gocmdb.Reportable) (error) - TestReporting()
	[X] legacyHandler(o gocmdb.Reportable) (error) - TestReporting()
	[X] usbCiCheckinV1(o gocmdb.Registerable) (error) - TestCheckinCheckout()
	[X] usbCiCheckoutV1(o gocmdb.Auditable) ([]byte, error) - TestCheckinCheckout()
	[X] audit(o gocmdb.Auditable) (error) - TestAuditing()
	[X] usbCiAuditV1(o gocmdb.Auditable) (error) - TestAuditing()
	[X] readFile(string, []byte) (error) - TestFileReadWrite()
	[X] writeFile([]byte, string) (error) - TestFileReadWrite()

TODO (NEW):
	[ ] audit(dev usb.Auditer) (err error)
	[X] report(dev usb.Reporter) (err error)
	[X] serial(dev usb.Serializer) (err error)
	[X] usbCiNewSnV1(dev usb.Serializer) (string, error)
	[ ] usbCiCheckinV1(dev usb.Reporter) (error)
	[ ] usbCiCheckoutV1(dev usb.Auditer) ([]byte, error)
	[ ] usbCiAuditV1(dev usb.Auditer) (error)
	[ ] usbMetaVendorV1(dev usb.Updater) (s string, err error)
	[ ] usbMetaProductV1(dev usb.Updater) (s string, err error)
	[ ] httpPost(url string, j []byte ) (b []byte, hs httpStatus, err error)
	[ ] httpGet(url string) (b []byte, hs httpStatus, err error)
	[ ] httpRequest(req *http.Request) (b []byte, hs httpStatus, err error)
	[ ] newConfig(cf string) (*Config, error)
	[ ] loadConfig(t interface{}, cf string) error
	[ ] makePath(path string) (string, error)
	[ ] displayVersion()
	[ ] route(i interface{}) (err error)
	[ ] convert(i interface{}) (interface{}, error)
	[ ] update(i interface{}) (interface{})
*/
