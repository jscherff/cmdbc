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
	`flag`
	`log`
	`os`
	`path/filepath`
	`strings`

	`github.com/google/gousb`
	`github.com/jscherff/gocmdb/usbci`
)

var (
	conf *Config
	slog, clog, elog *log.Logger
)

func init() {

	var err error

	// Build systemwide configuration from config file.

	if conf, err = newConfig(`config.json`); err != nil {
		log.Fatalf(err.Error())
	}

	// Return if in testing mode.

	if conf.Testing { return }

	// Initialized loggers.

	slog, clog, elog = newLoggers()

	// If run as legacy app executable, skip flag processing.

	if strings.Contains(filepath.Base(os.Args[0]), `magtek_inventory`) {
		return
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

	var fs *flag.FlagSet

	switch {

	case *fActionReport:
		fs = fsReport

	case *fActionSerial:
		fs = fsSerial

	case *fActionAudit:
		fs = fsAudit
	}

	if fs.Parse(os.Args[2:]); fs.NFlag() == 0 {
		fmt.Fprintln(os.Stderr, `You must specify an option.`)
		fs.Usage()
		os.Exit(1)
	}
}

func main() {

	// Return if in testing mode.

	if conf.Testing { return }

	// Instantiate context to enumerate attached USB devices.

	ctx := gousb.NewContext()
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

	// Open devices that match selection criteria in the Include.ProductID
	// and Include.VendorID maps from the configuration file.

	devs, _ := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {

		vid, pid := desc.Vendor.String(), desc.Product.String()

		if val, ok := conf.Include.ProductID[vid][pid]; ok {
			return val
		}
		if val, ok := conf.Include.VendorID[vid]; ok {
			return val
		}

		return conf.Include.Default
	})

	// Log and exit if no relevant devices found.

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
