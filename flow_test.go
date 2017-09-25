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
	`path/filepath`
	`reflect`
	`strings`
	`testing`
	`github.com/jscherff/gocmdb`
	`github.com/jscherff/gocmdb/usbci`
	`github.com/jscherff/gotest`
)

var (
	mag3, mag4 *usbci.Magtek
	gen3, gen4 *usbci.Generic
)

func init() {

	var errM3, errM4, errG3, errG4 error

	if mag3, errM3 = usbci.NewMagtek(nil); errM3 == nil {
		errM3 = mag3.RestoreJSON(mag1JSON)
	}

	if mag4, errM4 = usbci.NewMagtek(nil); errM4 == nil {
		errM4 = mag4.RestoreJSON(mag2JSON)
	}

	if gen3, errG3 = usbci.NewGeneric(nil); errG3 == nil {
		errG3 = gen3.RestoreJSON(gen1JSON)
	}

	if gen4, errG4 = usbci.NewGeneric(nil); errG4 == nil {
		errG4 = gen4.RestoreJSON(gen2JSON)
	}

	if errM3 != nil || errM4 != nil || errG3 != nil || errG4 != nil {
		log.Fatal(`Testing setup failed: could not restore devices.`)
	}
}


/* TODO: test each path through the application by setting flags and
   passing an object to the router for both magtek and generic devices.

	Actions:

	  -audit
	        Audit devices
	  -checkin
	        Check devices in
	  -legacy
	        Legacy operation
	  -report
	        Report actions
	  -reset
	        Reset device
	  -serial
	        Set serial number


	Audit Options:

	  -local
	        Audit against local state
	  -server
	        Audit against server state


	Report Options:

	  -console
	        Write reports to console
	  -folder <path>
	        Write reports to <path>
	  -format <format>
	        Report <format> {csv|nvp|xml|json}


	Serial Options:

	  -copy
	        Copy factory serial number
	  -erase
	        Erase current serial number
	  -fetch
	        Fetch serial number from server
	  -force
	        Force serial number change
	  -set "24F9999"
	        Set serial number to <string>
*/

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

func restoreState(tb testing.TB, o gocmdb.Auditable, b []byte) {

	tb.Helper()

	var err error

	if err = o.RestoreJSON(b); err != nil {
		tb.Fatal(`Testing failed: could not restore devices to original state.`)
	}

	if err = checkinDevice(o); err != nil {
		tb.Fatal(`Testing failed: could not restore database to original state.`)
	}
}

func TestAuditFlow(t *testing.T) {

	t.Run(`Flags: -audit -local`, func(t *testing.T) {

		var err error

		// Reset and Set flags.

		resetFlags(t)
		*fActionAudit = true
		*fAuditLocal = true

		// Send device to router.

		err = magtekRouter(mag3)
		gotest.Assert(t, err != nil, `first run should result in file-not-found error`)

		// Determine whether there are no changes recorded when auditing same device.

		err = magtekRouter(mag3)
		gotest.Ok(t, err)

		gotest.Assert(t, len(mag3.Changes) == 0, `device change log should be empty`)

		// Determine whether device differences are recorded in device change log.

		err = magtekRouter(mag4)
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(mag4.Changes, magChanges),
			`device change log does not contain known device differences`)

		// Determine whether device differences are recorded in app change log.

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, ClogCh1) && strings.Contains(fs, ClogCh2),
			`application change log does not contain known device differences`)

		// Restore state

		restoreState(t, mag4, mag2JSON)
	})

	t.Run(`Flags: -audit -server`, func(t *testing.T) {

		var err error

		// Reset and Set flags.

		resetFlags(t)
		*fActionAudit = true
		*fAuditServer = true

		// Send device to router.

		err = magtekRouter(mag3)
		gotest.Ok(t, err)

		// Determine whether there are no changes recorded when auditing same device.

		err = magtekRouter(mag3)
		gotest.Ok(t, err)

		// Determine whether device differences are recorded in device change log.

		err = magtekRouter(mag4)
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(mag4.Changes, magChanges),
			`device change log does not contain known device differences`)

		// Determine whether device differences are recorded in app change log.

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, ClogCh1) && strings.Contains(fs, ClogCh2),
			`application change log does not contain known device differences`)

		// Restore state

		restoreState(t, mag4, mag2JSON)
	})
}

func TestCheckinFlow(t *testing.T) {

	t.Run(`Flags: -checkin`, func(t *testing.T) {

		var err error

		// Reset and Set flags.

		resetFlags(t)
		*fActionCheckin = true

		// Change a property.

		mag4.VendorName = `Check-in Test`

		// Send device to router.

		err = magtekRouter(mag4)
		gotest.Ok(t, err)

		// Checkout device and test if property change persisted.

		b, err := checkoutDevice(mag4)
		gotest.Ok(t, err)

		err = mag4.RestoreJSON(b)
		gotest.Ok(t, err)

		gotest.Assert(t, mag4.VendorName == `Check-in Test`, `device changes did not persist after checkin`)

		// Restore state

		restoreState(t, mag4, mag2JSON)
	})
}

func TestLegacyFlow(t *testing.T) {

	t.Run(`Flags: -legacy`, func(t *testing.T) {

		var err error

		// Reset and Set flags.

		resetFlags(t)
		*fActionLegacy = true

		// Send device to router.

		err = magtekRouter(mag1)
		gotest.Ok(t, err)

		// Test whether signature of report file content is correct.

		b, err := ioutil.ReadFile(conf.Files.Legacy)
		gotest.Ok(t, err)

		gotest.Assert(t, mag1SigLegacy == sha256.Sum256(b), `unexpected hash signature of Legacy report`)
	})
}

func TestReportFlow(t *testing.T) {

	t.Run(`Flags: -report -folder -format csv`, func(t *testing.T) {

		var err error

		// Reset and Set flags.

		resetFlags(t)
		*fActionReport = true
		*fReportFormat = `csv`

		// Send device to router.

		err = magtekRouter(mag1)
		gotest.Ok(t, err)

		// Test whether signature of report file content is correct.

		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, mag1SigCSV == sha256.Sum256(b), `unexpected hash signature of CSV report`)
	})

	t.Run(`Flags: -report -folder -format nvp `, func(t *testing.T) {

		var err error

		// Reset and Set flags.

		resetFlags(t)
		*fActionReport = true
		*fReportFormat = `nvp`

		// Send device to router.

		err = magtekRouter(mag1)
		gotest.Ok(t, err)

		// Test whether signature of report file content is correct.

		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, mag1SigNVP == sha256.Sum256(b), `unexpected hash signature of NVP report`)
	})

	t.Run(`Flags: -report -folder -format xml`, func(t *testing.T) {

		var err error

		// Reset and Set flags.

		resetFlags(t)
		*fActionReport = true
		*fReportFormat = `xml`

		// Send device to router.

		err = magtekRouter(mag1)
		gotest.Ok(t, err)

		// Test whether signature of report file content is correct.

		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, mag1SigPXML == sha256.Sum256(b), `unexpected hash signature of XML report`)
	})

	t.Run(`Flags: -report -folder -format json`, func(t *testing.T) {

		var err error

		// Reset and Set flags.

		resetFlags(t)
		*fActionReport = true
		*fReportFormat = `json`

		// Send device to router.

		err = magtekRouter(mag1)
		gotest.Ok(t, err)

		// Test whether signature of report file content is correct.

		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)
		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, mag1SigPJSON == sha256.Sum256(b), `unexpected hash signature of JSON report`)
	})
}

/*
	t.Run(`Flags: -serial -erase`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialErase = true

		ctx := gousb.NewContext()
		defer ctx.Close()

		dev, err := ctx.OpenDeviceWithVIDPID(0x0801, 0x0001)

		if err != nil {
			t.Skip(`no compatible devices found`)
		}

		defer dev.Close()


		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -copy`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialCopy = true
		*fSerialFetch = false
		*fSerialSet = ``

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -fetch`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialCopy = false
		*fSerialFetch = true
		*fSerialSet = ``

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -set "24F9999"`, func(t *testing.T) {
		*fActionSerial = false
		*fSerialCopy = false
		*fSerialFetch = false
		*fSerialSet = `24F9999`

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -force -copy`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialForce = true
		*fSerialCopy = true
		*fSerialFetch = false
		*fSerialSet = ``

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -force -fetch`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialForce = true
		*fSerialCopy = false
		*fSerialFetch = true
		*fSerialSet = ``

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -force -set "24F9999"`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialForce = true
		*fSerialCopy = false
		*fSerialFetch = false
		*fSerialSet = `24F9999`

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -erase -copy`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialErase = true
		*fSerialCopy = true
		*fSerialFetch = false
		*fSerialSet = ``

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -erase -fetch`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialErase = true
		*fSerialCopy = false
		*fSerialFetch = true
		*fSerialSet = ``

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -serial -erase -set "24F9999"`, func(t *testing.T) {
		*fActionSerial = true
		*fSerialErase = true
		*fSerialCopy = false
		*fSerialFetch = false
		*fSerialSet = `24F9999`

		magtekRouter(mag3)
		genericRouter(gen3)
	})

	t.Run(`Flags: -reset`, func(t *testing.T) {
		*fActionReset = true

		magtekRouter(mag3)
		genericRouter(gen3)
	})
*/

