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

// legacyHandler writes legacy report to application directory.
func legacyHandler(o gocmdb.Reportable) (err error) {
	if err = writeFile(o.Legacy(), conf.Files.Legacy); err != nil {
		elog.Print(err)
	}
	return err
}

// reportHandler processes report options and writes report to the
// selected destination.
func reportHandler(o gocmdb.Reportable) (err error) {

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

	f := fmt.Sprintf(`%s.%s`, o.Filename(), *fReportFormat)

	switch {

	case *fReportConsole:
		fmt.Fprintf(os.Stdout, string(b))

	case len(*fReportFolder) > 0:
		err = writeFile(b, filepath.Join(*fReportFolder, f))

	default:
		err = writeFile(b, filepath.Join(conf.Paths.ReportDir, f))
	}

	if err != nil {
		elog.Print(err)
	}

	return err
}

// serialHandler processes the serial number options and configures the
// the serial number.
func serialHandler(o gocmdb.Configurable) (err error) {

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

	if !*fSerialForce && o.ID() != `` {

		err = fmt.Errorf(`device %s-%s serial: SN already set to %q`,
			o.VID(), o.PID(), o.ID(),
		)

		elog.Print(err)
		return err
	}

	switch {

	case *fSerialSet != ``:

		slog.Printf(`device %s-%s serial: setting SN to %q`,
			o.VID(), o.PID(), *fSerialSet,
		)

		err = o.SetDeviceSN(*fSerialSet)

	case *fSerialCopy:

		slog.Printf(`device %s-%s serial: copying factory SN`,
			o.VID(), o.PID(),
		)

		err = o.CopyFactorySN(7)

	case *fSerialFetch:

		if s, err = usbCiNewSnV1(o); err != nil {
			break
		}

		slog.Printf(`device %s-%s serial: setting SN %q from server`,
			o.VID(), o.PID(), s,
		)

		if err = o.SetDeviceSN(s); err != nil {
			break
		}

		slog.Printf(`device %s-%s-%s serial: checking in with server`,
			o.VID(), o.PID(), o.ID(),
		)

		err = usbCiCheckinV1(o)
	}

	if err != nil {
		elog.Print(err)
	}

	return err
}

// auditAdtion performs a change audit against a previously-safed local
// state file (-local) or against properties from the last server checkin.
// It then logs the changes to the local change log and reports changes
// back to the server.
func auditHandler(o gocmdb.Auditable) (err error) {

	var ch [][]string

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

		ch, err = o.CompareFile(f)

		slog.Printf(`device %s-%s-%s audit: saving current state to %q`,
			o.VID(), o.PID(), o.ID(), f,
		)

		if err := o.Save(f); err != nil {
			elog.Print(err) // Local scope.
		}

	case *fAuditServer:

		slog.Printf(`device %s-%s-%s audit: fetching previous state from server`,
			o.VID(), o.PID(), o.ID(),
		)

		var j []byte

		if j, err = usbCiCheckoutV1(o); err == nil {
			ch, err = o.CompareJSON(j)
		}

		slog.Printf(`device %s-%s-%s audit: saving current state to server`,
			o.VID(), o.PID(), o.ID(),
		)

		if err := usbCiCheckinV1(o); err != nil {
			elog.Print(err)
		}

	default:

		err = fmt.Errorf(`device %s-%s-%s audit: invalid audit option`,
			o.VID(), o.PID(), o.ID(),
		)
	}

	if err != nil {
		elog.Print(err)
		return err
	}

	if len(ch) == 0 {
		slog.Printf(`device %s-%s-%s audit: no changes`, o.VID(), o.PID(), o.ID())
		return nil
	}

	slog.Printf(`device %s-%s-%s audit: recording changes in change log`,
		o.VID(), o.PID(), o.ID(),
	)

	for _, c := range ch {
		clog.Printf(`device %s-%s-%s changed: %q was %q, now %q`,
			o.VID(), o.PID(), o.ID(), c[0], c[1], c[2],
		)
	}

	slog.Printf(`device %s-%s-%s audit: reporting changes to server`,
		o.VID(), o.PID(), o.ID(),
	)

	o.SetChanges(ch)
	err = usbCiAuditV1(o)

	if err != nil {
		elog.Print(err)
	}

	return err
}
