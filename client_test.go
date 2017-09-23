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

var (
	gen1JSON = []byte(
	`{
		"host_name": "John-SurfacePro",
		"vendor_id": "0acd",
		"product_id": "2030",
		"vendor_name": "ID TECH",
		"product_name": "TM3 Magstripe USB-HID Keyboard Reader",
		"serial_num": "",
		"software_id": "",
		"product_ver": "",
		"bus_number": 1,
		"bus_address": 14,
		"port_number": 1,
		"buffer_size": 0,
		"max_pkt_size": 8,
		"usb_spec": "2.00",
		"usb_class": "per-interface",
		"usb_subclass": "per-interface",
		"usb_protocol": "0",
		"device_speed": "full",
		"device_ver": "1.00",
		"object_type": "*usbci.Generic",
		"device_sn": "",
		"factory_sn": "",
		"descriptor_sn": "",
		"changes": null,
		"vendor": {}
	}`)

	gen2JSON = []byte(
	`{
		"host_name": "John-SurfacePro",
		"vendor_id": "0acd",
		"product_id": "2030",
		"vendor_name": "ID TECH",
		"product_name": "TM3 Magstripe USB-HID Keyboard Reader",
		"serial_num": "",
		"software_id": "",
		"product_ver": "",
		"bus_number": 1,
		"bus_address": 14,
		"port_number": 1,
		"buffer_size": 0,
		"max_pkt_size": 8,
		"usb_spec": "2.00",
		"usb_class": "per-interface",
		"usb_subclass": "per-interface",
		"usb_protocol": "0",
		"device_speed": "full",
		"device_ver": "1.00",
		"object_type": "*usbci.Generic",
		"device_sn": "",
		"factory_sn": "",
		"descriptor_sn": "",
		"changes": null,
		"vendor": {}
	}`)

	mag1JSON = []byte(
	`{
		"host_name": "John-SurfacePro",
		"vendor_id": "0801",
		"product_id": "0001",
		"vendor_name": "Mag-Tek",
		"product_name": "USB Swipe Reader",
		"serial_num": "24F0014",
		"software_id": "21042818B01",
		"product_ver": "",
		"bus_number": 1,
		"bus_address": 13,
		"port_number": 1,
		"buffer_size": 24,
		"max_pkt_size": 8,
		"usb_spec": "1.10",
		"usb_class": "per-interface",
		"usb_subclass": "per-interface",
		"usb_protocol": "0",
		"device_speed": "full",
		"device_ver": "1.00",
		"object_type": "*usbci.Magtek",
		"device_sn": "24F0014",
		"factory_sn": "",
		"descriptor_sn": "24F0014",
		"changes": null,
		"vendor": {}
	}`)

	mag2JSON = []byte(
	`{
		"host_name": "John-SurfacePro",
		"vendor_id": "0801",
		"product_id": "0001",
		"vendor_name": "Mag-Tek",
		"product_name": "USB Swipe Reader",
		"serial_num": "24F0014",
		"software_id": "21042818B02",
		"product_ver": "",
		"bus_number": 1,
		"bus_address": 13,
		"port_number": 1,
		"buffer_size": 24,
		"max_pkt_size": 8,
		"usb_spec": "2.00",
		"usb_class": "per-interface",
		"usb_subclass": "per-interface",
		"usb_protocol": "0",
		"device_speed": "full",
		"device_ver": "1.00",
		"object_type": "*usbci.Magtek",
		"device_sn": "24F0014",
		"factory_sn": "",
		"descriptor_sn": "24F0014",
		"changes": null,
		"vendor": {}
	}`)

	mag1ShaJSON = [32]byte{
		0xaf,0x81,0xad,0x6f,0xf7,0x6c,0x37,0xbd,
		0x45,0x8c,0xe8,0xfc,0xa5,0xd9,0x06,0x38,
		0x5b,0xc2,0x80,0x32,0x08,0x53,0x8b,0xac,
		0x86,0xe0,0x2a,0xdd,0xc9,0x8a,0x7a,0x32,
	}

	mag1ShaPrettyJSON = [32]byte{
		0x36,0x54,0xc7,0x2f,0x3e,0xf5,0xe3,0x4d,
		0xc8,0x67,0x66,0x17,0x27,0x9d,0x0e,0x1a,
		0xc0,0xde,0x50,0x0d,0x20,0x8e,0x54,0x33,
		0x00,0x9e,0x17,0x32,0xe1,0x90,0x0a,0xe7,
	}

	mag1ShaXML = [32]byte{
		0x82,0xc7,0x14,0x84,0xee,0xb5,0x4a,0x91,
		0xfc,0x92,0xa6,0x8b,0xeb,0xf7,0xd4,0x66,
		0x93,0xad,0xc0,0x6b,0x89,0x3e,0x99,0x11,
		0x28,0xfc,0x7e,0x61,0xf3,0x4f,0x7c,0xed,
	}

	mag1ShaPrettyXML = [32]byte{
		0x0f,0x05,0x4e,0x13,0x51,0x5e,0x90,0x9d,
		0x3d,0x39,0xfb,0xb8,0x6a,0x14,0x20,0xcb,
		0x3a,0xd0,0xb6,0x79,0xa5,0x56,0xad,0xf7,
		0xce,0xff,0x31,0xdc,0x56,0x2a,0xbd,0x92,
	}

	mag1ShaCSV = [32]byte{
		0x98,0xd5,0xe9,0x1d,0x6f,0xa9,0xe8,0xfe,
		0x7c,0xd6,0xa8,0xa0,0x7e,0x88,0x48,0xd4,
		0xcf,0x8b,0x04,0x9c,0x05,0x3e,0x1b,0x58,
		0x41,0x3c,0xf8,0x3e,0x27,0x8a,0x98,0xea,
	}

	mag1ShaNVP = [32]byte{
		0xd0,0xc4,0xea,0x8b,0x3c,0x80,0xae,0x79,
		0xe8,0x0e,0x17,0x1e,0xd3,0x55,0x09,0x88,
		0xbb,0x2b,0x11,0x84,0xac,0x3d,0xd9,0x42,
		0x50,0xc4,0x5d,0x5e,0x70,0xd3,0x65,0xe2,
	}

	mag1ShaLegacy = [32]byte{
		0xb3,0xb5,0x58,0x2b,0xb2,0xd9,0x88,0x4a,
		0x78,0xd5,0xf4,0x2d,0x98,0x0c,0x2b,0x81,
		0xfd,0xd1,0x43,0xb6,0xcc,0x58,0x14,0x39,
		0x23,0x30,0x50,0x2f,0xe3,0x59,0x88,0x5a,
	}

	mag1, mag2 *usbci.Magtek
	gen1, gen2 *usbci.Generic

	magChanges = make([][]string, 2)
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

func TestgetNewSN(t *testing.T) {

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
		gotest.Assert(t, mag1ShaPrettyJSON == sha256.Sum256(b), `unexpected hash signature of JSON report`)
	})

	t.Run("XML Report", func(t *testing.T) {

		*fReportFormat = `xml`
		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err := reportAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1ShaPrettyXML == sha256.Sum256(b), `unexpected hash signature of XML report`)
	})

	t.Run("CSV Report", func(t *testing.T) {

		*fReportFormat = `csv`
		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err := reportAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1ShaCSV == sha256.Sum256(b), `unexpected hash signature of CSV report`)
	})

	t.Run("NVP Report", func(t *testing.T) {

		*fReportFormat = `nvp`
		fn := filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.` + *fReportFormat)

		err := reportAction(mag1)
		gotest.Ok(t, err)

		b, err := ioutil.ReadFile(fn)
		gotest.Ok(t, err)
		gotest.Assert(t, mag1ShaNVP == sha256.Sum256(b), `unexpected hash signature of NVP report`)
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
		//[[SoftwareID 21042818B02 21042818B01] [USBSpec 2.00 1.10]]
		//[[SoftwareID 21042818B02 21042818B01] [USBSpec 2.00 1.10]]

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
log.Println(mag1mod.Changes)
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
log.Println(mag1)
log.Println(mag1mod)
		err = auditAction(mag1mod)
		gotest.Ok(t, err)
log.Println(mag1mod.Changes)
		gotest.Assert(t, reflect.DeepEqual(mag1mod.Changes, magChanges),
			`device change log does not contain known device differences`)

		fb, err := ioutil.ReadFile(conf.Files.ChangeLog)
		gotest.Ok(t, err)

		fs := string(fb)
		gotest.Assert(t, strings.Contains(fs, ch1) && strings.Contains(fs, ch2),
			`application change log does not contain known device differences`)
	})
}
