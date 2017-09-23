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

	return err // Errors already logged.
}

// serialAction processes the serial number options and configures the
// the serial number.
func serialAction(o gocmdb.Configurable) (err error) {

	var s string

	if *fSerialErase {

		slog.Printf(`device %s-%s serial: erasing SN %q`,
			o.VID(), o.PID(), o.ID(),
		)

		if err = o.EraseDeviceSN(); err != nil {
			elog.Print(err)
			return err
		}
	}

	if len(o.ID()) > 0 && !*fSerialForce {

		err = fmt.Errorf(`device %s-%s serial: SN already set to %q`,
			o.VID(), o.PID(), o.ID(),
		)

		elog.Print(err)
		return err
	}

	switch {

	case len(*fSerialSet) > 0:

		slog.Printf(`device %s-%s serial: setting SN to %q`,
			o.VID(), o.PID(), *fSerialSet,
		)

		if err = o.SetDeviceSN(*fSerialSet); err != nil {
			elog.Print(err)
		}

	case *fSerialCopy:

		slog.Printf(`device %s-%s serial: copying factory SN`,
			o.VID(), o.PID(),
		)

		if err = o.CopyFactorySN(7); err != nil {
			elog.Print(err)
		}

	case *fSerialFetch:

		if s, err = getNewSN(o); err != nil {
			break // Errors already logged.
		}

		slog.Printf(`device %s-%s serial: setting SN %q from server`,
			o.VID(), o.PID(), s,
		)

		if err = o.SetDeviceSN(s); err != nil {
			elog.Print(err)
			break
		}

		slog.Printf(`device %s-%s-%s serial: checking in with server`,
			o.VID(), o.PID(), o.ID(),
		)

		err = checkinDevice(o) // Errors already logged.
	}

	return err
}

// auditAdtion performs a change audit against a previously-safed local
// state file (-local) or against properties from the last server checkin.
// It then logs the changes to the local change log and reports changes
// back to the server.
func auditAction(o gocmdb.Auditable) (err error) {

	var chgs [][]string

	if o.ID() == `` {
		slog.Printf(`device %s-%s audit: skipping, no SN`, o.VID(), o.PID())
		return err
	}

	switch true {

	case *fAuditLocal:

		f := filepath.Join(
			conf.Paths.StateDir,
			fmt.Sprintf(`%s-%s-%s.json`, o.VID(), o.PID(), o.ID()),
		)

		slog.Printf(`device %s-%s-%s audit: fetching previous state from %q`,
			o.VID(), o.PID(), o.ID(), f,
		)

		if chgs, err = o.CompareFile(f); err != nil {
			elog.Print(err)
		}

		slog.Printf(`device %s-%s-%s audit: saving current state to %q`,
			o.VID(), o.PID(), o.ID(), f,
		)

		if errSave := o.Save(f); errSave != nil {
			elog.Print(errSave)
		}

	case *fAuditServer:

		slog.Printf(`device %s-%s-%s audit: fetching previous state from server`,
			o.VID(), o.PID(), o.ID(),
		)

		var c []byte

		if c, err = checkoutDevice(o); err == nil {
			if chgs, err = o.CompareJSON(c); err != nil {
				elog.Print(err)
			}
		}

		slog.Printf(`device %s-%s-%s audit: saving current state to server`,
			o.VID(), o.PID(), o.ID(),
		)

		checkinDevice(o) // Errors already logged.

	default:

		err = fmt.Errorf(`device %s-%s-%s audit: invalid audit option`,
			o.VID(), o.PID(), o.ID(),
		)

		elog.Print(err)
	}

	if err != nil {
		return err
	}

	if len(chgs) == 0 {

		slog.Printf(`device %s-%s-%s audit: no changes`, o.VID(), o.PID(), o.ID())

	} else {

		slog.Printf(`device %s-%s-%s audit: recording changes in change log`,
			o.VID(), o.PID(), o.ID(),
		)

		for _, chg := range chgs {
			clog.Printf(`device %s-%s-%s changed: %q was %q, now %q`,
				o.VID(), o.PID(), o.ID(),
				chg[0], chg[1], chg[2],
			)
		}

		slog.Printf(`device %s-%s-%s audit: reporting changes to server`,
			o.VID(), o.PID(), o.ID(),
		)

		o.SetChanges(chgs)

		err = submitAudit(o) // Errors already logged.
	}

	return err
}
