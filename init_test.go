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
	`flag`
	`fmt`
	`log`
	`os`
	`sync`
	`testing`
	`github.com/google/gousb`
	`github.com/jscherff/cmdb/ci/peripheral/usb`
)

type TestData struct {
	Jsn map[string][]byte
	Mag map[string]*usb.Magtek
	Gen map[string]*usb.Generic
	Idt map[string]*usb.IDTech
	Sig map[string]map[string][32]byte
	Chg [][]string
	Clg []string
}

var (
	td *TestData
	mux sync.Mutex
	testConfFile = `config.json`
	testDataFile = `tdata.json`
)

func init() {

	var err error

	if conf, err = newConfig(testConfFile); err != nil {
		log.Fatal(err)
	}

	td = &TestData{}

	if err = loadConfig(td, testDataFile); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func resetFlags(tb testing.TB) {

	tb.Helper()

	*fActionAudit = false
	*fActionCheckin = false
	*fActionReport = false
	*fActionReset = false
	*fActionSerial = false
	*fActionVersion = false

	*fReportConsole = false
	*fReportFolder = conf.Paths.ReportDir
	*fReportFormat = ``

	*fSerialDefault = false
	*fSerialErase = false
	*fSerialForce = false
	*fSerialFetch = false
	*fSerialSet = ``
}

func restoreState(tb testing.TB) {

	tb.Helper()

	if err := loadConfig(td, testDataFile); err != nil {
		tb.Fatal(err)
	}
}

func getMagtekDevice(tb testing.TB, c *gousb.Context) (*usb.Magtek, error) {

	tb.Helper()

	if dev, _ := c.OpenDeviceWithVIDPID(0x0801, 0x0001); dev != nil {
		return usb.NewMagtek(dev)
	} else {
		return nil, fmt.Errorf(`device not found`)
	}
}

