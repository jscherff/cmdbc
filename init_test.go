package main

import (
	`flag`
	`log`
	`os`
	`sync`
	`testing`
	`github.com/google/gousb`
	`github.com/jscherff/gocmdb/usbci`
)

var (
	mag1, mag2 *usbci.Magtek
	gen1, gen2 *usbci.Generic

	magChanges = make([][]string, 2)

	ClogCh1 = `"SoftwareID" was "21042840G01", now "21042840G02"`
	ClogCh2 = `"USBSpec" was "1.10", now "2.00"`

	mux sync.Mutex
)

func init() {
	magChanges[0] = []string{`SoftwareID`, `21042840G01`, `21042840G02`}
	magChanges[1] = []string{`USBSpec`, `1.10`, `2.00`}
}

func TestMain(m *testing.M) {

	var err error

	flag.Parse()

	if conf, err = newConfig(`config.json`); err != nil {
		log.Fatalf(err.Error())
	}

	conf.Logging.System.Console = false
	conf.Logging.Change.Console = false
	conf.Logging.Error.Console = false
	slog, clog, elog = newLoggers()

	if err = createObjects(); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func createObjects() (err error) {

	if mag1, err = usbci.NewMagtek(nil); err != nil {
		return err
	}

	if mag2, err = usbci.NewMagtek(nil); err != nil {
		return err
	}

	if gen1, err = usbci.NewGeneric(nil); err != nil {
		return err
	}

	if gen2, err = usbci.NewGeneric(nil); err != nil {
		return err
	}

	if err = mag1.RestoreJSON(mag1JSON); err != nil {
		return err
	}

	if err = mag2.RestoreJSON(mag2JSON); err != nil {
		return err
	}

	if err = gen1.RestoreJSON(gen1JSON); err != nil {
		return err
	}

	if err = gen2.RestoreJSON(gen2JSON); err != nil {
		return err
	}

	return err
}

func resetFlags(tb testing.TB) {

	tb.Helper()

	*fActionAudit = false
	*fActionCheckin = false
	*fActionLegacy = false
	*fActionReport = false
	*fActionReset = false
	*fActionSerial = false

	*fReportFolder = conf.Paths.ReportDir
	*fReportConsole = false
	*fReportFormat = ``

	*fSerialCopy = false
	*fSerialErase = false
	*fSerialForce = false
	*fSerialFetch = false
	*fSerialSet = ``

	*fAuditLocal = false
	*fAuditServer = false
}

func restoreState(tb testing.TB) {

	tb.Helper()

	if err := createObjects(); err != nil {
		tb.Fatal(err)
	}
}

func getMagtekDevice(tb testing.TB, c *gousb.Context) (mdev *usbci.Magtek, err error) {

	tb.Helper()

	dev, err := c.OpenDeviceWithVIDPID(0x0801, 0x0001)

	if dev != nil {
		mdev, err = usbci.NewMagtek(dev)
	}

	return mdev, err
}

