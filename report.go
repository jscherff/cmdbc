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
	`fmt`
	`io/ioutil`
	`os`
	`path/filepath`
	`github.com/jscherff/cmdb/ci/peripheral/usb`
)

// report processes options and writes report to the selected destination.
func report(dev usb.Reporter) (err error) {

	if fsReport.Parse(os.Args[2:]); fsReport.NFlag() == 0 {
		fsReport.Usage()
		os.Exit(1)
	}

	var b []byte

	switch *fReportFormat {

	case `csv`:
		b, err = dev.CSV()

	case `nvp`:
		b, err = dev.NVP()

	case `xml`:
		b, err = dev.PrettyXML()

	case `json`:
		b, err = dev.PrettyJSON()

	default:
		err = fmt.Errorf(`invalid format %q`, *fReportFormat)
	}

	if err != nil {
		return err
	}

	if *fReportConsole {
		fmt.Fprintln(os.Stdout, string(b))
		return nil
	}

	f := fmt.Sprintf(`%s-%s.%s`, dev.SN(), dev.Conn(), *fReportFormat)

	if *fReportFolder != `` {
		f = filepath.Join(*fReportFolder, f)
	} else {
		f = filepath.Join(conf.Paths.ReportDir, f)
	}

	return ioutil.WriteFile(f, b, FileMode)
}
