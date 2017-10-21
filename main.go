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
	`log`
	`os`
	`github.com/google/gousb`
)

const defaultConfig = `config.json`

var (
	conf *Config
	slog, clog, elog *log.Logger
)

func main() {

	var err error
        log.SetFlags(log.Flags() | log.Lshortfile)

	// Process command-line flags.

	if len(os.Args) < 2 {
		fsAction.Usage()
		os.Exit(1)
	}

	fsAction.Parse(os.Args[1:2])

	if *fActionVersion {
                displayVersion()
                os.Exit(0)
        }

	// Build system-wide configuration from config file.

	if conf, err = newConfig(defaultConfig); err != nil {
		log.Fatal(err)
	}

	// Initialize loggers.

	slog, clog, elog = newLoggers()

	// Instantiate context to enumerate devices.

	ctx := gousb.NewContext()
	ctx.Debug(conf.DebugLevel)
	defer ctx.Close()

	// Open devices that match selection criteria.

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {

		vid, pid := desc.Vendor.String(), desc.Product.String()

		if val, ok := conf.Include.ProductID[vid][pid]; ok {
			return val
		}
		if val, ok := conf.Include.VendorID[vid]; ok {
			return val
		}

		return conf.Include.Default
	})

	if err != nil && conf.DebugLevel > 0 {
		elog.Print(err)
	}

	// Exit if no devices found.

	if len(devs) == 0 {
		elog.Fatalf(`no devices found`)
	}

	// Pass each device to router.

	for _, dev := range devs {

		defer dev.Close()

		slog.Printf(`found device %s-%s`, dev.Desc.Vendor, dev.Desc.Product)

		if err = route(dev); err != nil {
			elog.Print(err)
		}
	}
}
