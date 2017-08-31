// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	 http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"os"
)

// Process report action and options.
func reportAction(o gocmdb.Reportable) (e error) {

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
		e = fmt.Errorf("report: invalid format %q", *fReportFormat)
	}

	if e == nil {

		switch {

		case len(*fReportFile) > 0:
			d, f := filepath.Split(*fReportFile)
			if len(d) == 0 {d = config.ReportDir}
			e = writeFile(b, filepath.Join(d, f))

		case *fReportStdout:
			fmt.Fprintf(os.Stdout, string(b))

		default:
			e = fmt.Errorf("report: no destintion")
		}

	}

	return e
}

// Process serial number action and options.
func serialAction(o gocmdb.Configurable, i gocmdb.Registerable) (e error) {

	if *fSerialErase {
		e = o.EraseDeviceSN()
	}

	s, e := o.DeviceSN()

	if len(s) != 0 && !*fSerialForce {
		e = fmt.Errorf("serial: already configured")
	}

	if e == nil {

		switch {

		case len(*fSerialConfig) > 0:
			e = o.SetDeviceSN(*fSerialConfig)

		case *fSerialCopy:
			e = o.CopyFactorySN(7)

		case *fSerialServer:
			s, e := serialRequest(i)
			if e != nil {
				break
			}
			if len(s) == 0 {
				e = fmt.Errorf("serial: empty serial number")
				break
			}
			e = o.SetDeviceSN(s)
		}
	}

	return e
}

// Process reset action.
func resetAction(o gocmdb.Resettable) (error) {
	return o.Reset()
}
