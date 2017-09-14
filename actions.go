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
	"fmt"
	"os"
	"path/filepath"
	"github.com/jscherff/gocmdb"
	`github.com/jscherff/goutils`
)

// legacyAction writes legacy report to application directory.
func legacyAction(o gocmdb.Reportable) (err error) {

	err = writeFile(o.Legacy(), filepath.Join(conf.Paths.AppDir, conf.Files.Legacy))

	if err != nil {
		err = gocmdb.ErrorDecorator(err)
	}

	return err
}

// reportAction processes report options and writes report to the
// selected destination.
func reportAction(o gocmdb.Reportable) (err error) {

	var b []byte

	switch *fReportFormat {

	case "csv":
		b, err = o.CSV()

	case "nvp":
		b, err = o.NVP()

	case "xml":
		b, err = o.PrettyXML()

	case "json":
		b, err = o.PrettyJSON()

	default:
		err = fmt.Errorf("invalid format %q", *fReportFormat)
	}

	if err == nil {

		switch {

		case *fReportConsole:
			fmt.Fprintf(os.Stdout, string(b))

		case len(*fReportFolder) > 0:
			err = writeFile(b, filepath.Join(*fReportFolder, o.Filename()))

		default:
			f := fmt.Sprintf("%s.%s", o.Filename(), *fReportFormat)
			err = writeFile(b, filepath.Join(conf.Paths.ReportDir, f))
		}

	}

	if err != nil {
		err = gocmdb.ErrorDecorator(err)
	}

	return err
}

// serialAction processes the serial number options and configures the
// the serial number.
func serialAction(o gocmdb.Configurable) (err error) {

	var s string

	if *fSerialErase {
		err = o.EraseDeviceSN()
	}

	if err == nil {
		s = o.ID()

		if len(s) > 0 && !*fSerialForce {
			err = fmt.Errorf("serial number already set to %q", s)
		}
	}

	if err == nil {

		switch {

		case len(*fSerialSet) > 0:
			err = o.SetDeviceSN(*fSerialSet)

		case *fSerialCopy:
			err = o.CopyFactorySN(7)

		case *fSerialFetch:
			if s, err = fetchSnRequest(o); err != nil {
				break
			}
			if len(s) == 0 {
				err = fmt.Errorf("empty serial number from server")
				break
			}
			if err = o.SetDeviceSN(s); err != nil {
				err = checkinRequest(o)
			}
		}
	}

	if err != nil {
		err = gocmdb.ErrorDecorator(err)
	}

	return err
}
