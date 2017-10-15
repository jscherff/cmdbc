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
	`strings`

	`github.com/google/gousb`
	`github.com/jscherff/gocmdb/usbci`
)

const defaultConfig = `config.json`

var (
	conf *Config
	slog, clog, elog *log.Logger
)

func main() {

	var err error

	// Build system-wide configuration from config file.

	if conf, err = newConfig(defaultConfig); err != nil {
		log.Fatalf(err.Error())
	}

	// Initialize loggers.

	slog, clog, elog = newLoggers()

	// Instantiate context to enumerate devices.

	ctx := gousb.NewContext()
	ctx.Debug(conf.DebugLevel)
	defer ctx.Close()

	// If run as legacy app executable, find first device matching magtek
	// vendor ID and product ID, produce legacy report, then exit.

	if strings.Contains(filepath.Base(os.Args[0]), `magtek_inventory`) {

		dev, err := ctx.OpenDeviceWithVIDPID(
			gousb.ID(usbci.MagtekVID),
			gousb.ID(usbci.MagtekPID),
		)

		if err != nil {
			elog.Fatal(err)
		}

		mdev, err := usbci.NewMagtek(dev)

		if err != nil {
			elog.Fatal(err)
		}

		legacyAction(mdev)

		os.Exit(0)
	}

	// Process command-line actions and options.

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, `You must specify an action.`)
		fsAction.Usage()
		os.Exit(1)
	}

	// Parse action flag.

	fsAction.Parse(os.Args[1:2])

	// Parse option flags associated with selected action flag.

	switch {

	case *fActionReport:
		if fsReport.Parse(os.Args[2:]); fsReport.NFlag() == 0 {
			fmt.Fprintln(os.Stderr, `You must specify an option.`)
			fsReport.Usage()
			os.Exit(1)
		}

	case *fActionSerial:
		if fsSerial.Parse(os.Args[2:]); fsSerial.NFlag() == 0 {
			fmt.Fprintln(os.Stderr, `You must specify an option.`)
			fsSerial.Usage()
			os.Exit(1)
		}

	case *fActionAudit:
		if fsAudit.Parse(os.Args[2:]); fsAudit.NFlag() == 0 {
			fmt.Fprintln(os.Stderr, `You must specify an option.`)
			fsAudit.Usage()
			os.Exit(1)
		}
	}

	// Open devices that match selection criteria in the Include.ProductID
	// and Include.VendorID maps from the configuration file.

	var openFunc = func(desc *gousb.DeviceDesc) bool {
fmt.Printf("%#v\n", desc)
		vid, pid := desc.Vendor.String(), desc.Product.String()

		if val, ok := conf.Include.ProductID[vid][pid]; ok {
			return val
		}
		if val, ok := conf.Include.VendorID[vid]; ok {
			return val
		}

		return conf.Include.Default
	}

	devs, err := ctx.OpenDevices(openFunc)

for _, d := range devs {
	fmt.Println(d.Desc.Bus, d.Desc.Address, d.Desc.Port)
}

	// Log and exit if no relevant devices found.
	if err != nil && conf.DebugLevel > 0 {
		elog.Print(err)
	}
	if len(devs) == 0 {
		elog.Fatalf(`no devices found`)
	}

	// Pass devices to relevant device handlers.

	for _, dev := range devs {

		defer dev.Close()

		slog.Printf(`found USB device: VID %s PID %s`,
			dev.Desc.Vendor.String(),
			dev.Desc.Product.String(),
		)

		switch uint16(dev.Desc.Vendor) {

		case usbci.MagtekVID:

			if d, err := usbci.NewMagtek(dev); err != nil {
				elog.Print(err)
			} else {
				slog.Printf(`identified USB device as %s`, d.Type())
				magtekRouter(d)
			}

		default:

			if d, err := usbci.NewGeneric(dev); err != nil {
				elog.Print(err)
			} else {
				slog.Printf(`identified USB device as %s`, d.Type())
				genericRouter(d)
			}
		}
	}
}
