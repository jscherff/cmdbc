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
	`fmt`
	`log`
	`math`
	`github.com/jscherff/gocmdb/usbci`
)

var (
	mag1JSON = []byte(
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
	)

	mag2JSON = []byte(
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
	)

	gen1JSON = []byte(
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
	)

	gen2JSON = []byte(
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
	)

	mag1, mag2 *usbci.Magtek
	gen1, gen2 *usbci.Generic
)

func printb(b [32]byte) (s string) {

	for i, b := range b {
		if math.Mod(float64(i), 8) == 0 { s += "\n\t\t" }
		s += fmt.Sprintf("0x%02x,", b)
	}

	return s
}

func init() {

	var errM1, errM2, errG1, errG2 error

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
		log.Fatalln("Testing setup failed: could not restore devices.")
	}
}

func main() {

	var (
		b []byte
		err error
	)

	fmt.Printf("package main\n\nvar (\n")

	fmt.Printf("\tmag1JSON = []byte(\n\t\t`%s`,\n\t)\n\n", mag1JSON)
	fmt.Printf("\tmag2JSON = []byte(\n\t\t`%s`,\n\t)\n\n", mag2JSON)
	fmt.Printf("\tgen1JSON = []byte(\n\t\t`%s`,\n\t)\n\n", gen1JSON)
	fmt.Printf("\tgen2JSON = []byte(\n\t\t`%s`,\n\t)\n\n", gen2JSON)

	if b, err = mag1.PrettyJSON(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag1SigPJSON = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag1.PrettyXML(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag1SigPXML = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag1.JSON(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag1SigJSON = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag1.XML(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag1SigXML = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag1.CSV(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag1SigCSV = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag1.NVP(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag1SigNVP = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b = mag1.Legacy(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag1SigLegacy = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag2.PrettyJSON(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag2SigPJSON = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag2.PrettyXML(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag2SigPXML = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag2.JSON(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag2SigJSON = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag2.XML(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag2SigXML = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag2.CSV(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag2SigCSV = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = mag2.NVP(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag2SigNVP = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b = mag2.Legacy(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tmag2SigLegacy = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen1.PrettyJSON(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen1SigPJSON = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen1.PrettyXML(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen1SigPXML = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen1.JSON(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen1SigJSON = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen1.XML(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen1SigXML = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen1.CSV(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen1SigCSV = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen1.NVP(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen1SigNVP = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b = gen1.Legacy(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen1SigLegacy = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen2.PrettyJSON(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen2SigPJSON = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen2.PrettyXML(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen2SigPXML = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen2.JSON(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen2SigJSON = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen2.XML(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen2SigXML = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen2.CSV(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen2SigCSV = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b, err = gen2.NVP(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen2SigNVP = [32]byte{%s\n\t}\n\n", printb(sha256.Sum256(b)))

	if b = gen2.Legacy(); err != nil { log.Fatalln(err) }
	fmt.Printf("\tgen2SigLegacy = [32]byte{%s\n\t}\n", printb(sha256.Sum256(b)))

	fmt.Printf(")\n\n")
}

