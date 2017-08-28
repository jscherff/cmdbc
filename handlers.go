package main

import (
	"github.com/jscherff/gocmdb/usbci/magtek"
	//"bytes"
	"fmt"
	"os"
)

func report(d *magtek.Device) (e error) {

	//var s string
	var b []byte

	di, errs := magtek.NewDeviceInfo(d)

	if len(errs) > 0 {
		e = fmt.Errorf("Error(s) getting device information")
	} else {

		switch *fReportFormat {

		case "csv":
			//r, e = di.CSV(!*fReportAll)
			b, e = di.CSV(!*fReportAll)

		case "nvp":
			//r, e = di.NVP(!*fReportAll)
			b, e = di.NVP(!*fReportAll)

		case "xml":
			//if e == nil {r = string(b)}
			b, e = di.XML(!*fReportAll)

		case "json":
			//if e == nil {r = string(b)}
			b, e = di.JSON(!*fReportAll)

		case "leg":
			b = []byte(fmt.Sprintf("%s,%s\n", di.HostName, di.SerialNum))

		default:
			e = fmt.Errorf("invalid report format %q", *fReportFormat)
		}
	}

	if e == nil {

		switch {

		case len(*fReportFile) > 0:
			//TODO

		case len(*fReportServer) > 0:
			//TODO

		case *fReportStdout:
			//fmt.Fprintf(os.Stdout, string(b))
			fmt.Fprintf(os.Stdout, string(b))

		default:
			e = fmt.Errorf("no report destintion selected")
		}

	}

	return e
}

func config(d *magtek.Device) (e error) {

	s, e := d.GetDeviceSN()

	if e == nil {

		switch {

		case *fConfigErase:
			e = d.EraseDeviceSN()
			fallthrough

		case len(s) > 0 && !*fConfigForce:
			e = fmt.Errorf("serial number already configured")

		case *fConfigCopy:
			e = d.CopyFactorySN(7)

		case len(*fConfigString) > 0:
			e = d.SetDeviceSN(*fConfigString)

		case len(*fConfigServer) > 0:
			e = d.SetDeviceSN("24F0000") //TODO: call server

		default:
			e = fmt.Errorf("nothing to do")
		}
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
