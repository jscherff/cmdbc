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
)

// legacyAction writes legacy report to application directory.
func legacyAction(o gocmdb.Reportable) (err error) {

	err = writeFile(o.Legacy(), filepath.Join(conf.Paths.AppDir, conf.Files.Legacy))

	if err != nil {
		elog.Print(err)
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
		elog.Print(err)
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
			elog.Print(err)
			return err
		}
	}

	if len(o.ID()) > 0 && !*fSerialForce {
		err = fmt.Errorf(`serial number already set to %q`, o.ID())
		elog.Print(err)
		return err
	}

	switch {

	case len(*fSerialSet) > 0:
		err = o.SetDeviceSN(*fSerialSet)
		elog.Print(err)

	case *fSerialCopy:
		err = o.CopyFactorySN(7)
		elog.Print(err)

	case *fSerialFetch:

		if s, err = newSNRequest(o); err != nil {
			// Error already decorated and logged.
			break
		}

		if len(s) == 0 {
			err = fmt.Errorf(`empty serial number from server`)
			elog.Print(err)
			break
		}

		if err = o.SetDeviceSN(s); err != nil {
			elog.Print(err)
			break
		}

		if err = checkinRequest(o); err != nil {
			elog.Print(err)
		}
	}

	return err
}

// auditAdtion requests a server-side audit against the previous checkin.
func auditAction(o gocmdb.Auditable) (err error) {

	var chgs [][]string
	b := 

	if o.ID() == `` {
		slog.Print(`skipping audit for VID %q PID %q: no serial number`, o.VID(), o.PID())
		return err
	}

	f := filepath.Join(conf.Paths.StateDir, fmt.Sprintf(`%s-%s-%s.json`, o.VID(), o.PID(), o.ID()))
	fi, err := os.Stat(f)

	if err == nil {
		slog.Printf(`found state file %q last modified %s`, fi.Name(), fi.ModTime())
		chgs, err = o.CompareFile(f)
	}

	if sverr := o.Save(f); sverr != nil {
		elog.Print(sverr)
	}

	if err != nil {
		elog.Print(err)
		return err
	}

	if len(chgs) > 0 {
		for _, chg := range chgs {
			clog.Printf(`device %s-%s-%s last audited %s: %q was %q, now %q`,
				o.VID(), o.PID(), o.ID(), fi.ModTime(), chg[0], chg[1], chg[2])
		}
	}

	//if j, err := json.Marshal(chgs); err == nil {

	// TODO: report to server
	// o.Changes = chgs

	return err
}
