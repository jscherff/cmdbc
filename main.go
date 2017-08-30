package main

import (
	"github.com/jscherff/gocmdb/usbci/magtek"
	"github.com/jscherff/gocmdb/usbci"
	"github.com/google/gousb"
	"errors"
	"log"
	"fmt"
	"os"
)

var conf *Config

func init() {

	var err error

	conf, err = GetConfig("config.json")

	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
}

func main() {

	var err error

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "You must specify an action.\n")
		fsAction.Usage()
		os.Exit(1)
	}

	fsAction.Parse(os.Args[1:2])

	switch {

	case *fActionReport:
		if fsReport.Parse(os.Args[2:]); fsReport.NFlag() == 0 {
			fsReport.Usage()
			os.Exit(1)
		}

	case *fActionConfig:
		if fsConfig.Parse(os.Args[2:]); fsReport.NFlag() == 0 {
			fsConfig.Usage()
			os.Exit(1)
		}
	}

	context := gousb.NewContext()
	defer context.Close()

	// Open devices that report a Magtek vendor ID, 0x0801.
	// We omit error checking on OpenDevices() because this
	// function terminates with 'libusb: not found [code -5]'
	// on Windows systems.

	devices, _ := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {

		vid := desc.Vendor.String()
		pid := desc.Product.String()

		if val, ok := conf.IncludePID[vid][pid]; ok {return val}
		if val, ok := conf.IncludeVID[vid]; ok {return val}

		return conf.IncludeDefault
	})

	if len(devices) == 0 {
		log.Fatalf("no devices found")
	}

	for _, device := range devices {

		defer device.Close()

		switch uint16(device.Desc.Vendor) {

		case magtek.MagtekVendorID:

			var mdev *magtek.Device
			var info *magtek.DeviceInfo

			mdev, err = magtek.NewDevice(device)

			if err == nil {
				info, err = magtek.NewDeviceInfo(mdev)
			}

			if err == nil {
				switch {

				case *fActionReport:
					err = report(info)

				case *fActionConfig:
					err = config(mdev)

				case *fActionReset:
					err = reset(mdev)

				default:
					err = errors.New("action not supported")
				}
			}

		default:

			var gdev *usbci.Device
			var info *usbci.DeviceInfo
			gdev, err = usbci.NewDevice(device)

			if err == nil {
				info, err = usbci.NewDeviceInfo(gdev)
			}

			if err == nil {
				switch {

				case *fActionReport:
					err = report(info)

				case *fActionReset:
					err = reset(gdev)

				default:
					err = errors.New("action not supported")
				}
			}
		}

		if err != nil {
			log.Printf("%v", err)
		}
	}
}
