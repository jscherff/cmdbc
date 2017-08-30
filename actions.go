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
			e = fmt.Errorf("no report destintion")
		}

	}

	return e
}

func serial(o gocmdb.Configurable) (e error) {

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

		case *fConfigServer:
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

func audit(o gocmdb.Reportable) (error) {
	return nil
}

func checkin(o gocmdb.Reportable) (error) {
	return nil
}

func writeFile(s string, b []byte) (e error) {

	d, f := filepath.Split(s)

	if len(d) == 0 {
		d = config.AppDir
	}

	p := fmt.Sprintf("%s%c%s", d, filepath.Separator, f)

	if e = os.MkdirAll(d, 0755); e == nil {
		e = ioutil.WriteFile(p, b, 0644)
	}

	return e
}
