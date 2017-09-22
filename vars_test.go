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
	`github.com/jscherff/gocmdb/usbci`
)

var (
	genJSNv1 = []byte(
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

	genJSNv2 = []byte(
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

	magJSNv1 = []byte(
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

	magJSNv2 = []byte(
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

	mag1ShaJSN = [32]byte{
		0x36,0x54,0xc7,0x2f,0x3e,0xf5,0xe3,0x4d,
		0xc8,0x67,0x66,0x17,0x27,0x9d,0x0e,0x1a,
		0xc0,0xde,0x50,0x0d,0x20,0x8e,0x54,0x33,
		0x00,0x9e,0x17,0x32,0xe1,0x90,0x0a,0xe7,
	}

	mag1ShaXML = [32]byte{
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

	mag1, mag2 *usbci.Magtek
	gen1, gen2 *usbci.Generic
)
