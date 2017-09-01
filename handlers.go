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
	"errors"

	"github.com/google/gousb"
	"github.com/jscherff/gocmdb/usbci"
	"github.com/jscherff/gocmdb/usbci/magtek"
)

func handleMagtek(d *gousb.Device) (e error) {

	var md *magtek.Device
	var mdi *magtek.DeviceInfo

	md, e = magtek.NewDevice(d)

	if e == nil {
		mdi, e = magtek.NewDeviceInfo(md)
	}

	if e == nil {
		switch {

		case *fActionReport:
			e = reportAction(mdi)

		case *fActionSerial:
			e = serialAction(md, mdi)
			defer resetAction(md)

		case *fActionReset:
			e = resetAction(md)

		case *fActionAudit:
			e = auditRequest(mdi)

		case *fActionCheckin:
			e = checkinRequest(mdi)

		default:
			e = errors.New("action not supported")
		}
	}

	return e
}

func handleGeneric(d *gousb.Device) (e error) {

	var gd *usbci.Device
	var gdi *usbci.DeviceInfo

	gd, e = usbci.NewDevice(d)

	if e == nil {
		gdi, e = usbci.NewDeviceInfo(gd)
	}

	if e == nil {
		switch {

		case *fActionReport:
			e = reportAction(gdi)

		case *fActionReset:
			e = resetAction(gd)

		case *fActionAudit:
			e = auditRequest(gdi)

		case *fActionCheckin:
			e = checkinRequest(gdi)

		default:
			e = errors.New("action not supported")
		}
	}

	return e
}
