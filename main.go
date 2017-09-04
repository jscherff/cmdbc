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
	"fmt"
	"log"
	"os"

	"github.com/jscherff/gocmdb/usbci/magtek"
	"github.com/google/gousb"
)

// The config variable holds the runtime configuration.
var config *Config

func init() {

	var e error
	config, e = getConfig()

	if e != nil {
		log.Fatalf("error processing config: %v", e)
	}
}

func main() {

	var e error

	// Process command-line actions and options.

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "You must specify an action.\n")
		fsAction.Usage()
		os.Exit(1)
	}

	// Parse action flag.

	fsAction.Parse(os.Args[1:2])

	// Parse option flags associated with selected action flag.

	switch {

	case *fActionReport:
		if fsReport.Parse(os.Args[2:]); fsReport.NFlag() == 0 {
			fmt.Fprintf(os.Stderr, "You must specify an option.\n")
			fsReport.Usage()
			os.Exit(1)
		}

	case *fActionSerial:
		if fsSerial.Parse(os.Args[2:]); fsSerial.NFlag() == 0 {
			fmt.Fprintf(os.Stderr, "You must specify an option.\n")
			fsSerial.Usage()
			os.Exit(1)
		}
	}

	// Instantiate context to enumerate attached USB devices.

	context := gousb.NewContext()
	defer context.Close()

	// Open devices that match selection criteria in the IncludePID
	// and IncludeVID maps from the configuration file.

	devices, _ := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {

		vid, pid := desc.Vendor.String(), desc.Product.String()

		if val, ok := config.IncludePID[vid][pid]; ok {
			return val
		}
		if val, ok := config.IncludeVID[vid]; ok {
			return val
		}

		return config.DefaultInclude
	})

	// Log and exit if no relevant devices found.

	if len(devices) == 0 {
		log.Fatalf("no devices found")
	}

	// Pass devices to relevant device handlers.

	for _, device := range devices {

		defer device.Close()

		switch uint16(device.Desc.Vendor) {

		case magtek.MagtekVendorID:
			e = magtekHandler(device)

		default:
			e = genericHandler(device)
		}

		if e != nil {
			log.Printf("%v", e)
		}
	}
}
