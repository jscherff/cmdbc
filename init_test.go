package main

import (
	`crypto/sha256`
	`flag`
	`log`
	`os`
	`testing`
	`github.com/google/gousb`
	`github.com/jscherff/gocmdb/usbci`
)

func init() {

	magChanges[0] = []string{`SoftwareID`, `21042840G01`, `21042840G02`}
	magChanges[1] = []string{`USBSpec`, `1.10`, `2.00`}

	if err := createObjects(); err != nil {
		log.Fatal(err)
	}

	if err := generateSigs(); err != nil {
		log.Fatal(err)
	}
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

	for k, j := range magJSON {

		if mag[k], err = usbci.NewMagtek(nil); err != nil {
			return err
		}
		if err := mag[k].RestoreJSON(j); err != nil {
			return err
		}
	}

	for k, j := range genJSON {

		if genDev[k], err = usbci.NewGeneric(nil); err != nil {
			return err
		}
		if err := genDev[k].RestoreJSON(j); err != nil {
			return err
		}
	}

	return nil
}

func generateSigs() error {

	for k := range magJSON {

		if b, err := mag[k].CSV(); err != nil {
			return err
		} else {
			sigCSV[k] = sha256.Sum256(b)
		}
		if b, err := mag[k].NVP(); err != nil {
			return err
		} else {
			sigNVP[k] = sha256.Sum256(b)
		}
		if b, err := mag[k].XML(); err != nil {
			return err
		} else {
			sigXML[k] = sha256.Sum256(b)
		}
		if b, err := mag[k].JSON(); err != nil {
			return err
		} else {
			sigJSON[k] = sha256.Sum256(b)
		}
		if b, err := mag[k].PrettyXML(); err != nil {
			return err
		} else {
			sigPrettyXML[k] = sha256.Sum256(b)
		}
		if b, err := mag[k].PrettyJSON(); err != nil {
			return err
		} else {
			sigPrettyJSON[k] = sha256.Sum256(b)
		}

		b := mag[k].Legacy()
		sigLegacy[k] = sha256.Sum256(b)
	}

	for k := range genJSON {

		if b, err := genDev[k].CSV(); err != nil {
			return err
		} else {
			sigCSV[k] = sha256.Sum256(b)
		}
		if b, err := genDev[k].NVP(); err != nil {
			return err
		} else {
			sigNVP[k] = sha256.Sum256(b)
		}
		if b, err := genDev[k].XML(); err != nil {
			return err
		} else {
			sigXML[k] = sha256.Sum256(b)
		}
		if b, err := genDev[k].JSON(); err != nil {
			return err
		} else {
			sigJSON[k] = sha256.Sum256(b)
		}
		if b, err := genDev[k].PrettyXML(); err != nil {
			return err
		} else {
			sigPrettyXML[k] = sha256.Sum256(b)
		}
		if b, err := genDev[k].PrettyJSON(); err != nil {
			return err
		} else {
			sigPrettyJSON[k] = sha256.Sum256(b)
		}

		b := genDev[k].Legacy()
		sigLegacy[k] = sha256.Sum256(b)
	}

	return nil
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

