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
	`github.com/google/gousb`
	`github.com/jscherff/cmdb/ci/peripheral/usb`
)

var checkin = usbCiCheckinV1 // Alias

func convert(i interface{}) (d interface{}, err error) {

	var v, p gousb.ID

	switch t := i.(type) {

	case *gousb.Device:
		v = t.Desc.Vendor
		p = t.Desc.Product

	case *gousb.DeviceDesc:
		v = t.Vendor
		p = t.Product

	default:
		return nil, fmt.Errorf(`unsupported type %T`, t)
	}

	switch {

	case usb.IsMagtek(v, p):
		return usb.NewMagtek(i)

	case usb.IsIDTech(v, p):
		return usb.NewIDTech(i)

	default:
		return usb.NewGeneric(i)
	}
}

func route(i interface{}) (error) {

	i, err := convert(i)

	if err != nil {
		return err
	}

	if d, ok := i.(usb.Serializer); ok {

		switch {

		case *fActionSerial:
			if err = serial(d); err != nil {
				return err
			}
			*fActionReset = true
		}
	}

	if d, ok := i.(usb.Auditer); ok {

		switch {

		case *fActionReport:
			err = report(d)

		case *fActionCheckin:
			err = checkin(d)

		case *fActionAudit:
			err = audit(d)

		case *fActionReset:
			err = d.Reset()
		}
	}

	return err
}
