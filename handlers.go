package main

import (
	"github.com/jscherff/gocmdb"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"os"
)

func report(o gocmdb.Reportable) (e error) {

	var b []byte

	switch *fReportFormat {

	case "csv":
		b, e = o.CSV()

	case "nvp":
		b, e = o.NVP()

	case "xml":
		b, e = o.XML()

	case "json":
		b, e = o.JSON()

	case "bare":
		b = o.Bare()

	default:
		e = fmt.Errorf("invalid report format %q", *fReportFormat)
	}

	if e == nil {

		switch {

		case len(*fReportFile) > 0:
			e = writeFile(*fReportFile, b)

		case *fReportStdout:
			fmt.Fprintf(os.Stdout, string(b))

		default:
			e = fmt.Errorf("no report destintion selected")
		}

	}

	return e
}

func config(o gocmdb.Configurable) (e error) {

	if *fConfigErase {
		e = o.EraseDeviceSN()
	}

	s, e := o.DeviceSN()

	if e == nil {

		switch {

		case len(s) > 0 && !*fConfigForce:
			e = fmt.Errorf("serial number already configured")

		case len(*fConfigString) > 0:
			e = o.SetDeviceSN(*fConfigString)

		case len(*fConfigServer) > 0:
			e = o.SetDeviceSN("24F0000") //TODO: call server

		case *fConfigCopy:
			e = o.CopyFactorySN(7)

		default:
			e = fmt.Errorf("nothing to do")
		}
	}

	return e
}

func reset(o gocmdb.Resettable) (error) {
	return o.Reset()
}

func writeFile(s string, b []byte) (e error) {

	d, f := filepath.Split(s)

	if len(d) == 0 {
		d = conf.AppPath
	}

	p := fmt.Sprintf("%s%c%s", d, filepath.Separator, f)
	fmt.Println(p)

	if e = os.MkdirAll(d, 0755); e == nil {
		e = ioutil.WriteFile(p, b, 0644)
	}

	return e
}
