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
	`fmt`
	`log`
	`os`
	`path/filepath`
	`reflect`
	`testing`
	`github.com/google/gousb`
	`github.com/jscherff/gocmdb/usbci`
	`github.com/jscherff/gotest`
)

var (
	gen1JSN = []byte(
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

	gen2JSN = []byte(
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

	mag1JSN = []byte(
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

	mag2JSN = []byte(
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

	mag1, mag2 *usbci.Magtek
	gen1, gen2 *usbci.Generic

	magChanges = make([][]string, 2)
)

func init() {

	magChanges[0] = []string{`SoftwareID`, `21042818B01`, `21042818B02`}
	magChanges[1] = []string{`USBSpec`, `1.10`, `2.00`}

	var errM1, errM2, errG1, errG2 error

	if mag1, errM1 = usbci.NewMagtek(nil); errM1 == nil {
		errM1 = mag1.RestoreJSON(mag1JSN)
	}

	if mag2, errM2 = usbci.NewMagtek(nil); errM2 == nil {
		errM2 = mag2.RestoreJSON(mag2JSN)
	}

	if gen1, errG1 = usbci.NewGeneric(nil); errG1 == nil {
		errG1 = gen1.RestoreJSON(gen1JSN)
	}

	if gen2, errG2 = usbci.NewGeneric(nil); errG2 == nil {
		errG2 = gen2.RestoreJSON(gen2JSN)
	}

	if errM1 != nil || errM2 != nil || errG1 != nil || errG2 != nil {
		log.Fatal(os.Stderr, "Testing setup failed: could not restore devices.")
	}
}

func TestGetterMethods(t *testing.T) {

	gotest.Assert(t, mag1.ID() == mag1.SerialNum, `ID() does not match (device).SerialNum`)
	gotest.Assert(t, mag1.VID() == mag1.VendorID, `VID() does not match (device).VenndorID`)
	gotest.Assert(t, mag1.PID() == mag1.ProductID, `PID() does not match (device).ProductID`)
	gotest.Assert(t, mag1.Type() == reflect.TypeOf(mag1).String(), `Type does not match TypeOf(device)`)

	if hostName, err := os.Hostname(); err != nil {
		return
	} else {
		gotest.Assert(t, mag1.Host() == hostName, `Host() does not match os.Hostname()`)
	}
}

func TestFilenameMethod(t *testing.T) {

	fileName := fmt.Sprintf(`%03d-%03d-%03d-%s-%s`,
                mag1.BusNumber,
                mag1.BusAddress,
                mag1.PortNumber,
                mag1.VendorID,
                mag1.ProductID,
        )

	gotest.Assert(t, mag1.Filename() == fileName, `(device).Filename() string incorrect`)
}

func TestPersistenceMethods(t *testing.T) {

	var err error

	t.Run("Save() and RestoreFile()", func(t *testing.T) {

		fn := filepath.Join(os.Getenv(`TEMP`), `mag1.json`)

		err = mag1.Save(fn)
		gotest.Ok(t, err)

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreFile(fn)
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(mag1, mag3), `restored device not identical to saved device`)
	})

	t.Run("JSON() and RestoreJSON()", func(t *testing.T) {

		j, err := mag1.JSON()
		gotest.Ok(t, err)

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreJSON(j)
		gotest.Ok(t, err)

		gotest.Assert(t, reflect.DeepEqual(mag1, mag3), `restored device not identical to saved device`)
	})
}

func TestCompareMethods(t *testing.T) {

	t.Run("Save() and CompareFile()", func(t *testing.T) {

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreJSON(mag1JSN)
		gotest.Ok(t, err)

		fn1 := filepath.Join(os.Getenv(`TEMP`), `mag1.json`)
		fn2 := filepath.Join(os.Getenv(`TEMP`), `mag2.json`)

		err = mag1.Save(fn1)
		gotest.Ok(t, err)

		err = mag2.Save(fn2)
		gotest.Ok(t, err)

		ss1, err := mag3.CompareFile(fn1)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss1) == 0, `cloned device should match parent device`)

		ss2, err := mag3.CompareFile(fn2)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss2) == 2, `cloned device should not match a modified device`)
	})

	t.Run("JSON() and CompareJSON()", func(t *testing.T) {

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreJSON(mag1JSN)
		gotest.Ok(t, err)

		j1, err := mag1.JSON()
		gotest.Ok(t, err)

		j2, err := mag2.JSON()
		gotest.Ok(t, err)

		ss1, err := mag3.CompareJSON(j1)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss1) == 0, `cloned device should match parent device`)

		ss2, err := mag3.CompareJSON(j2)
		gotest.Ok(t, err)
		gotest.Assert(t, len(ss2) == 2, `cloned device should not match a modified device`)
	})
}

func TestAuditMethods(t *testing.T) {

	t.Run("Save() and AuditFile()", func(t *testing.T) {

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreJSON(mag1JSN)
		gotest.Ok(t, err)

		fn1 := filepath.Join(os.Getenv(`TEMP`), `mag1.json`)
		fn2 := filepath.Join(os.Getenv(`TEMP`), `mag2.json`)

		err = mag1.Save(fn1)
		gotest.Ok(t, err)

		err = mag2.Save(fn2)
		gotest.Ok(t, err)

		err = mag3.AuditFile(fn1)
		gotest.Ok(t, err)
		gotest.Assert(t, len(mag3.Changes) == 0, `cloned device should match parent device`)

		err = mag3.AuditFile(fn2)
		gotest.Ok(t, err)
		gotest.Assert(t, len(mag3.Changes) == 2, `cloned device should not match a modified device`)

		if len(mag3.Changes) < 2 { return }

		gotest.Assert(t, reflect.DeepEqual(mag3.Changes, magChanges),
			`(device).Changes contains bad data`)
		gotest.Assert(t, reflect.DeepEqual(mag3.GetChanges(), magChanges),
			`(device).GetChanges() returns bad data`)
	})

	t.Run("JSON() and AuditJSON()", func(t *testing.T) {

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreJSON(mag1JSN)
		gotest.Ok(t, err)

		j1, err := mag1.JSON()
		gotest.Ok(t, err)

		j2, err := mag2.JSON()
		gotest.Ok(t, err)

		err = mag3.AuditJSON(j1)
		gotest.Ok(t, err)
		gotest.Assert(t, len(mag3.Changes) == 0, `cloned device should match parent device`)

		err = mag3.AuditJSON(j2)
		gotest.Ok(t, err)
		gotest.Assert(t, len(mag3.Changes) == 2, `cloned device should not match a modified device`)

		if len(mag3.Changes) < 2 { return }

		gotest.Assert(t, reflect.DeepEqual(mag3.Changes, magChanges),
			`(device).Changes contains bad data`)
		gotest.Assert(t, reflect.DeepEqual(mag3.GetChanges(), magChanges),
			`(device).GetChanges() returns bad data`)
	})
}

func TestChangeMethods(t *testing.T) {

	t.Run("AddChange() and GetChanges()", func(t *testing.T) {

		var changes = []string{`SoftwareID`, `21042818B01`, `21042818B03`}

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreJSON(mag1JSN)
		gotest.Ok(t, err)

		mag3.AddChange(`SoftwareID`, `21042818B01`, `21042818B03`)
		gotest.Assert(t, len(mag3.Changes) == 1, `(device).Changes should contain one change`)
		gotest.Assert(t, len(mag3.GetChanges()) == 1, `(device).GetChanges() should contain one change`)

		if len(mag3.Changes) < 1 { return }

		gotest.Assert(t, reflect.DeepEqual(mag3.GetChanges()[0], changes),
			`(device).GetChanges() returns bad data`)
		gotest.Assert(t, reflect.DeepEqual(mag3.Changes[0], changes),
			`(device).Changes contains bad data`)
	})

	t.Run("SetChanges() adn GetChanges()", func(t *testing.T) {

		mag3, err := usbci.NewMagtek(nil)
		gotest.Ok(t, err)

		err = mag3.RestoreJSON(mag1JSN)
		gotest.Ok(t, err)

		ss, err := mag3.CompareJSON(mag2JSN)
		gotest.Ok(t, err)

		mag3.SetChanges(ss)
		gotest.Assert(t, len(mag3.Changes) == 2, `(device).Changes should contain two changes`)
		gotest.Assert(t, len(mag3.GetChanges()) == 2, `(device).GetChanges() should contain two changes`)

		if len(mag3.Changes) < 2 { return }

		gotest.Assert(t, reflect.DeepEqual(mag3.Changes, magChanges),
			`(device).Changes contains bad data`)
		gotest.Assert(t, reflect.DeepEqual(mag3.GetChanges(), magChanges),
			`(device).GetChanges() returns bad data`)
	})
}

func TestSerialMethods(t *testing.T) {

	t.Run("Magtek Sureswipe Card Reader", func(t *testing.T) {

		ctx := gousb.NewContext()
		defer ctx.Close()

		dev, err := ctx.OpenDeviceWithVIDPID(0x0801, 0x0001)
		gotest.Ok(t, err)

		defer dev.Close()

		// Set device SN

		mdev, err := usbci.NewMagtek(dev)
		gotest.Ok(t, err)

		oldSn, err := mdev.GetDeviceSN()
		gotest.Ok(t, err)

		err = mdev.SetDeviceSN(`TESTING`)
		gotest.Ok(t, err)

		newSn, err := mdev.GetDeviceSN()
		gotest.Ok(t, err)
		gotest.Assert(t, newSn == `TESTING`, `setting device SN to new value unsuccessful`)

		errs := mdev.Refresh()
		gotest.Assert(t, len(errs) == 0, `(device).Refresh() produced setter errors`)

		// Erase device SN

		err = mdev.EraseDeviceSN()
		gotest.Ok(t, err)

		newSn, err = mdev.GetDeviceSN()
		gotest.Ok(t, err)
		gotest.Assert(t, newSn == ``, `erasing device SN was unsuccessful`)

		// Restore device SN

		err = mdev.SetDeviceSN(oldSn)
		gotest.Ok(t, err)
		newSn, err = mdev.GetDeviceSN()
		gotest.Ok(t, err)
		gotest.Assert(t, newSn == oldSn, `restoring device SN to previous value unsuccessful`)

		err = mdev.Reset()
		gotest.Ok(t, err)
	})

	t.Run("Magtek Dynamag Card Reader", func(t *testing.T) {

		ctx := gousb.NewContext()
		defer ctx.Close()

		dev, err := ctx.OpenDeviceWithVIDPID(0x0801, 0x0001)
		gotest.Ok(t, err)

		defer dev.Close()

		mdev, err := usbci.NewMagtek(dev)
		gotest.Ok(t, err)
		gotest.Assert(t, mdev.FactorySN != ``, `device does not have a factory SN`)

		oldSn, err := mdev.GetDeviceSN()
		gotest.Ok(t, err)

		err = mdev.CopyFactorySN(7)
		gotest.Ok(t, err)
		gotest.Assert(t, mdev.DeviceSN == mdev.FactorySN[:7], `copying factory SN to device SN unsuccessful`)

		err = mdev.SetDeviceSN(oldSn)
		gotest.Ok(t, err)
		newSn, err := mdev.GetDeviceSN()
		gotest.Ok(t, err)
		gotest.Assert(t, newSn == oldSn, `restoring device SN to previous value unsuccessful`)

		err = mdev.Reset()
		gotest.Ok(t, err)
	})
}
