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
	`errors`
	`fmt`
	`os`
	`path/filepath`
	`github.com/jscherff/gocmdb`
)

// legacyAction writes legacy report to application directory.
func legacyAction(o gocmdb.Reportable) (err error) {

	err = WriteFile(o.Legacy(), filepath.Join(conf.Paths.AppDir, conf.Files.Legacy))

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
		err = WriteFile(b, filepath.Join(*fReportFolder, o.Filename()))

	default:
		f := fmt.Sprintf(`%s.%s`, o.Filename(), *fReportFormat)
		err = WriteFile(b, filepath.Join(conf.Paths.ReportDir, f))
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

		if s, err = GetNewSN(o); err != nil {
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

		if err = SubmitCheckin(o); err != nil {
			elog.Print(err)
		}
	}

	return err
}

// auditAdtion requests a server-side audit against the previous checkin.
func auditAction(o gocmdb.Auditable) (err error) {

	var chgs [][]string

	if o.ID() == `` {
		slog.Printf(`skipping audit for VID %q PID %q: no serial number`, o.VID(), o.PID())
		return err
	}

	switch true {

	case *fAuditLocal:

		f := filepath.Join(conf.Paths.StateDir, fmt.Sprintf(`%s-%s-%s.json`, o.VID(), o.PID(), o.ID()))
		fi, err := os.Stat(f)

		if err == nil {
			slog.Printf(`found state file %q last modified %s`, fi.Name(), fi.ModTime())
			chgs, err = o.CompareFile(f)
		}

		if sErr := o.Save(f); sErr != nil {
			elog.Print(sErr)
		}

	case *fAuditServer:

		c, err := GetDevice(o)

		if err == nil {
			chgs, err = o.CompareJSON(c)
		}

	default:

		err = errors.New(`invalid audit option`)
	}

	if err != nil {
		elog.Print(err)
		return err
	}

	if len(chgs) == 0 {

		slog.Printf(`device %s-%s-%s audited: no changes`, o.VID(), o.PID(), o.ID())

	} else {

		slog.Printf(`device %s-%s-%s audited: changes recorded in change log`,
			o.VID(), o.PID(), o.ID(),
		)

		for _, chg := range chgs {
			clog.Printf(`device %s-%s-%s changed: %q was %q, now %q`,
				o.VID(), o.PID(), o.ID(), chg[0], chg[1], chg[2],
			)
		}

		o.SetChanges(chgs)

		if err = SubmitAudit(o); err != nil {
			elog.Print(err)
		}
	}

	return err
}
