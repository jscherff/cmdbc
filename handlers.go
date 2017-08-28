package main

import (
	"github.com/jscherff/gocmdb/usbci/magtek"
	"fmt"
	"os"
)

func report(d *magtek.Device) (e error) {

	var r string
	di, errs := magtek.NewDeviceInfo(d)

	if len(errs) > 0 {
		e = fmt.Errorf("Errors encountered getting device information")
	} else {

		switch *fReportFormat {

		case "csv":
			r, e = di.CSV(!*fReportAll)

		case "nvp":
			r, e = di.NVP(!*fReportAll)

		case "xml":
			b, e := di.XML(!*fReportAll)
			if e == nil {r = string(b)}

		case "json":
			b, e := di.JSON(!*fReportAll)
			if e == nil {r = string(b)}

		default:
			fmt.Fprintf(os.Stderr, "Invalid report format.\n")
			fsReport.Usage()
			os.Exit(1)
		}
	}

	if e == nil {

		switch {

		case len(*fReportFile) > 0:
			//TODO

		case *fReportStdout:
			fmt.Fprintf(os.Stdout, r)

		default:
			fmt.Fprintf(os.Stderr, "No report destination selected.\n")
			fsReport.Usage()
			os.Exit(1)
		}

	}

	return e
}

func config(d *magtek.Device) (e error) {

	s, e := d.GetDeviceSN()
	if e != nil {return e}

	switch {

	case *fConfigErase:
		e = d.EraseDeviceSN()
		fallthrough

	case len(s) > 0 && !*fConfigForce:
		fmt.Fprintf(os.Stderr, "Serial number already set. Exiting.\n")

	case *fConfigCopy:
		e = d.CopyFactorySN(7)

	case len(*fConfigSet) > 0:
		e = d.SetDeviceSN(*fConfigSet)

	case len(*fConfigUrl) > 0:
		e = d.SetDeviceSN("24F0000") //TODO: call server

	default:
		//TODO
	}

	return e
}

func reset(d *magtek.Device) (e error) {

	switch {

	case *fResetUsb:
		e = d.Reset()

	case *fResetDev:
		e = d.DeviceReset()
	}

	return e
}

