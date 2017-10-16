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
	`encoding/json`
	`flag`
	`log`
	`os`
	`sync`
	`testing`
	`github.com/google/gousb`
	`github.com/jscherff/gocmdb/usbci`
)

type TestData struct {
	Mag map[string]*usbci.Magtek
	Gen map[string]*usbci.Generic
	Sig map[string]map[string][32]byte
	Chg [][]string
	Clg []string
}

var (
	td *TestData
	mux sync.Mutex
)

func init() {

	td = new(TestData)

	if err := createObjects(); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {

	var err error

	flag.Parse()

	if conf, err = newConfig(`config.json`); err != nil {
		log.Fatal(err)
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

func createObjects() error {

	if fh, err := os.Open(`testdata.json`); err != nil {
		return err
	} else {
		defer fh.Close()
		return json.NewDecoder(fh).Decode(&td)
	}
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

