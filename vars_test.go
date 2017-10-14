package main

import (
	`sync`

	`github.com/jscherff/gocmdb/usbci`
)

var (
	magJSON = map[string][]byte{

		`mag1`: []byte(
			`{
				"host_name": "John-SurfacePro",
				"vendor_id": "0801",
				"product_id": "0001",
				"serial_number": "24FFFFF",
				"vendor_name": "Mag-Tek",
				"product_name": "USB Swipe Reader",
				"product_ver": "V05",
				"firmware_ver": "",
				"software_id": "21042840G01",
				"bus_number": 1,
				"bus_address": 7,
				"port_number": 1,
				"buffer_size": 60,
				"max_pkt_size": 8,
				"usb_spec": "1.10",
				"usb_class": "per-interface",
				"usb_subclass": "per-interface",
				"usb_protocol": "0",
				"device_speed": "full",
				"device_ver": "1.00",
				"object_type": "*usbci.Magtek",
				"device_sn": "24FFFFF",
				"factory_sn": "B164F78022713AA",
				"descriptor_sn": "24FFFFF"
			}`,
		),

		`mag2`: []byte(
			`{
				"host_name": "John-SurfacePro",
				"vendor_id": "0801",
				"product_id": "0001",
				"serial_number": "24FFFFF",
				"vendor_name": "Mag-Tek",
				"product_name": "USB Swipe Reader",
				"product_ver": "V05",
				"firmware_ver": "",
				"software_id": "21042840G02",
				"bus_number": 1,
				"bus_address": 7,
				"port_number": 1,
				"buffer_size": 60,
				"max_pkt_size": 8,
				"usb_spec": "2.00",
				"usb_class": "per-interface",
				"usb_subclass": "per-interface",
				"usb_protocol": "0",
				"device_speed": "full",
				"device_ver": "1.00",
				"object_type": "*usbci.Magtek",
				"device_sn": "24FFFFF",
				"factory_sn": "B164F78022713AA",
				"descriptor_sn": "24FFFFF"
			}`,
		),
	}

	genJSON = map[string][]byte{

		`gen1`: []byte(
			`{
				"host_name": "John-SurfacePro",
				"vendor_id": "0acd",
				"product_id": "2030",
				"serial_number": "",
				"vendor_name": "ID TECH",
				"product_name": "TM3 Magstripe USB-HID Keyboard Reader",
				"product_ver": "",
				"firmware_ver": "",
				"software_id": "",
				"bus_number": 1,
				"bus_address": 8,
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
				"descriptor_sn": ""
			}`,
		),

		`gen2`: []byte(
			`{
				"host_name": "John-SurfacePro",
				"vendor_id": "0acd",
				"product_id": "2030",
				"serial_number": "",
				"vendor_name": "ID TECH",
				"product_name": "TM4 Magstripe USB-HID Keyboard Reader",
				"product_ver": "",
				"firmware_ver": "",
				"software_id": "",
				"bus_number": 1,
				"bus_address": 8,
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
				"descriptor_sn": ""
			}`,
		),
	}

	mag = make(map[string]*usbci.Magtek)
	genDev = make(map[string]*usbci.Generic)

	sigCSV = make(map[string][32]byte)
	sigNVP = make(map[string][32]byte)
	sigXML = make(map[string][32]byte)
	sigJSON = make(map[string][32]byte)
	sigPrettyXML = make(map[string][32]byte)
	sigPrettyJSON = make(map[string][32]byte)
	sigLegacy = make(map[string][32]byte)

	magChanges = make([][]string, 2)

	changeLogCh1 = `"SoftwareID" was "21042840G01", now "21042840G02"`
	changeLogCh2 = `"USBSpec" was "1.10", now "2.00"`

	mux sync.Mutex
)
