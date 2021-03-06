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

func route(i interface{}) (err error) {

	if i, err = convert(i); err != nil {
		return err
	}

	i = update(i)

	if d, ok := i.(usb.Serializer); ok {

		switch {

		case *fActionSerial:
			if err = serial(d); err != nil {
				return err
			}
			*fActionReset = true
		}
	}

	if d, ok := i.(usb.Analyzer); ok {

		switch {

		case *fActionState:
			if err = showState(d); err != nil {
				return err
			}
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

func convert(i interface{}) (interface{}, error) {

	var v, p gousb.ID

	switch t := i.(type) {

	case *gousb.Device:
		v = t.Desc.Vendor
		p = t.Desc.Product

	case *gousb.DeviceDesc:
		v = t.Vendor
		p = t.Product

	case *usb.Device:
		return t, nil

	case *usb.Generic:
		return t, nil

	case *usb.Magtek:
		return t, nil

	case *usb.IDTech:
		return t, nil

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

func update(i interface{}) (interface{}) {

	d, ok := i.(usb.Updater)

	if !ok {
		return i
	}

	if d.GetVendorName() == `` {
		if s, err := vendor(d); err == nil {
			d.SetVendorName(s)
		}
	}

	if d.GetProductName() == `` {
		if s, err := product(d); err == nil {
			d.SetProductName(s)
		}
	}

	return i
}
