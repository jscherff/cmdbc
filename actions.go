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
	//`encoding/json`
	`fmt`
	`os`
	`path/filepath`
	`github.com/jscherff/gocmdb`
	`github.com/jscherff/goutil`
)

// legacyAction writes legacy report to application directory.
func legacyAction(o gocmdb.Reportable) (err error) {

	err = writeFile(o.Legacy(), filepath.Join(conf.Paths.AppDir, conf.Files.Legacy))

	if err != nil {
		elog.Println(goutil.ErrorDecorator(err))
	}

	return err
}

// reportAction processes report options and writes report to the
// selected destination.
func reportAction(o gocmdb.Reportable) (err error) {

	var b []byte

	switch *fReportFormat {

	case `csv`:
		b, err = o.CSV()

	case `nvp`:
		b, err = o.NVP()

	case `xml`:
		b, err = o.PrettyXML()

	case `json`:
		b, err = o.PrettyJSON()

	default:
		err = fmt.Errorf(`invalid format %q`, *fReportFormat)
	}

	if err != nil {
		elog.Println(goutil.ErrorDecorator(err))
		return err
	}

	switch {

	case *fReportConsole:
		fmt.Fprintf(os.Stdout, string(b))

	case len(*fReportFolder) > 0:
		err = writeFile(b, filepath.Join(*fReportFolder, o.Filename()))

	default:
		f := fmt.Sprintf(`%s.%s`, o.Filename(), *fReportFormat)
		err = writeFile(b, filepath.Join(conf.Paths.ReportDir, f))
	}

	// Error already decorated and logged.
	return err
}

// serialAction processes the serial number options and configures the
// the serial number.
func serialAction(o gocmdb.Configurable) (err error) {

	var s string

	if *fSerialErase {
		if err = o.EraseDeviceSN(); err != nil {
			elog.Println(goutil.ErrorDecorator(err))
			return err
		}
	}

	if len(o.ID()) > 0 && !*fSerialForce {
		err = fmt.Errorf(`serial number already set to %q`, s)
		elog.Println(goutil.ErrorDecorator(err))
		return err
	}

	switch {

	case len(*fSerialSet) > 0:
		err = o.SetDeviceSN(*fSerialSet)
		elog.Println(goutil.ErrorDecorator(err))

	case *fSerialCopy:
		err = o.CopyFactorySN(7)
		elog.Println(goutil.ErrorDecorator(err))

	case *fSerialFetch:

		if s, err = fetchSnRequest(o); err != nil {
			// Error already decorated and logged.
			break
		}

		if len(s) == 0 {
			err = fmt.Errorf(`empty serial number from server`)
			elog.Println(goutil.ErrorDecorator(err))
			break
		}

		if err = o.SetDeviceSN(s); err != nil {
			elog.Println(goutil.ErrorDecorator(err))
			break
		}

		if err = checkinRequest(o); err != nil {
			elog.Println(goutil.ErrorDecorator(err))
		}
	}

	return err
}

// auditAdtion requests a server-side audit against the previous checkin.
func auditAction(o gocmdb.Auditable) (err error) {

	//var j []byte

	if o.ID() == `` {
		err = fmt.Errorf(`device with VID %q PID %q has no serial number`, o.VID(), o.PID())
		elog.Println(err.Error())
		return err
	}

	f := filepath.Join(conf.Paths.StateDir, fmt.Sprintf(`%s-%s-%s.json`, o.VID(), o.PID(), o.ID()))
	fi, err := os.Stat(f)

	if err != nil {
		elog.Println(err.Error())
		return err
	} else {
		slog.Printf(`found state file %q size %d last modified %s`, fi.Name(), fi.Size(), fi.ModTime())
	}

	chgs, err := o.CompareFile(f)

	if err != nil {
		elog.Println(err.Error())
		return err
	}

	if len(chgs) > 0 {
		for _, chg := range chgs {
			clog.Printf(`device %s-%s-%s since %s, property %q was %q, now %q`,
				o.VID(), o.PID(), o.ID(), fi.ModTime(), chg[0], chg[1], chg[2])
		}
	}

	// TODO: report to server
	// o.Changes = chgs

	if err = o.Save(f); err != nil {
		elog.Println(err.Error())
	}

	return err
}

/*
	fmt.Println("\nSaving 'test2.json'")
	o.Save("test2.json")

	fmt.Println("\no.CSV()")
	b, err := o.CSV()
	fmt.Println(string(b), err)

	fmt.Println("\no.JSON()")
	b, err = o.JSON()
	fmt.Println(string(b), err)

	fmt.Println("\no.XML()")
	b, err = o.XML()
	fmt.Println(string(b), err)

	fmt.Println("\no.NVP()")
	b, err = o.NVP()
	fmt.Println(string(b), err)

	fmt.Println("\nComparing to 'test.json'")
	ss, err := o.CompareFile("test.json")
	fmt.Println(ss)

	return err
*/
