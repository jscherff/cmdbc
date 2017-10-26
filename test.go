package main

var usbVendorData = map[string]string{
	`03f0`: `Hewlett-Packard`,
	`0467`: `AT&T Paradyne`,
	`0468`: `Wieson Technologies Co., Ltd`,
	`046a`: `Cherry GmbH`,
}

var usbProductData = map[string]map[string]string{
	`03f0`: map[string]string {
		`0004`: `DeskJet 895c`,
		`0011`: `OfficeJet G55`,
		`0012`: `DeskJet 1125C Printer Port`,
		`0024`: `KU-0316 Keyboard`,
		`002a`: `LaserJet P1102`,
		`0101`: `ScanJet 4100c`,
		`0102`: `PhotoSmart S20`,
		`0104`: `DeskJet 880c/970c`,
	},
	`046a`: map[string]string {
		`0001`: `Keyboard`,
		`0004`: `CyBoard Keyboard`,
		`0005`: `XX33 SmartCard Reader Keyboard`,
		`0008`: `Wireless Keyboard and Mouse`,
	},
}
