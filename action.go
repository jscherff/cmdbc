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

// audit performs a change audit against properties from the last checkin.
func audit(dev usb.Auditer) (err error) {

	var ch [][]string

	if dev.SN() == `` {
		sl.Printf(`device %s-%s skipping audit, no serial number`,
			dev.VID(), dev.PID(),
		)
		return err
	}

	sl.Printf(`device %s-%s-%s fetching previous state from server`,
		dev.VID(), dev.PID(), dev.SN(),
	)

	var j []byte

	if j, err = checkout(dev); err != nil {
		sl.Printf(`device %s-%s-%s skipping audit: no previous state`,
			dev.VID(), dev.PID(), dev.SN(),
		)
		return err
	}

	if ch, err = dev.CompareJSON(j); err != nil {
		return err
	}

	sl.Printf(`device %s-%s-%s saving current state to server`,
		dev.VID(), dev.PID(), dev.SN(),
	)

	if err := checkin(dev); err != nil {
		el.Print(err) // err occluded later by sendAudit()
	}

	if len(ch) == 0 {
		sl.Printf(`device %s-%s-%s detected no changes`,
			dev.VID(), dev.PID(), dev.SN(),
		)
		return nil
	}

	sl.Printf(`device %s-%s-%s recording changes in change log`,
		dev.VID(), dev.PID(), dev.SN(),
	)

	for _, c := range ch {
		cl.Printf(`device %s-%s-%s modified: '%s' was '%s', now '%s'`,
			dev.VID(), dev.PID(), dev.SN(), c[0], c[1], c[2],
		)
	}

	sl.Printf(`device %s-%s-%s reporting changes to server`,
		dev.VID(), dev.PID(), dev.SN(),
	)

	dev.SetChanges(ch)

	return sendAudit(dev)
}

// report processes options and writes report to the selected destination.
func report(dev usb.Reporter) (err error) {

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
		err = fmt.Errorf(`invalid format '%s'`, *fReportFormat)
	}

	if err != nil {
		return err
	}

	if *fReportConsole {
		fmt.Fprintln(os.Stdout, string(b))
		return nil
	}

	f := fmt.Sprintf(`%s-V%s-P%s.%s`, dev.Conn(), dev.VID(), dev.PID(), *fReportFormat)

	if *fReportFolder != `` {
		f = filepath.Join(*fReportFolder, f)
	} else {
		f = filepath.Join(conf.Paths.ReportDir, f)
	}

	return ioutil.WriteFile(f, b, FileMode)
}

// serial processes options and configures the the serial number.
func serial(dev usb.Serializer) (err error) {

	var s string

	if *fSerialErase {

		sl.Printf(`device %s-%s erasing serial '%s'`,
			dev.VID(), dev.PID(), dev.SN(),
		)

		if err = dev.EraseDeviceSN(); err != nil {
			return err
		}
	}

	if !*fSerialForce && dev.SN() != `` {

		return fmt.Errorf(`device %s-%s serial already set to '%s'`,
			dev.VID(), dev.PID(), dev.SN(),
		)

	}

	switch {

	case *fSerialFetch:

		if s, err = newSn(dev); err != nil {
			break
		}

		sl.Printf(`device %s-%s setting serial to '%s'`,
			dev.VID(), dev.PID(), s,
		)

		if err = dev.SetDeviceSN(s); err != nil {
			break
		}

		sl.Printf(`device %s-%s-%s checking in with server`,
			dev.VID(), dev.PID(), dev.SN(),
		)

		err = checkin(dev)

	case *fSerialDefault:

		sl.Printf(`device %s-%s setting serial to default`,
			dev.VID(), dev.PID(),
		)

		err = dev.SetDefaultSN()

	case *fSerialSet != ``:

		sl.Printf(`device %s-%s setting serial to '%s'`,
			dev.VID(), dev.PID(), *fSerialSet,
		)

		err = dev.SetDeviceSN(*fSerialSet)
	}

	return err
}

// showState displays any available diagnostic information about the device.
func showState(dev usb.Analyzer) (error) {

	if state, err := dev.GetState(); err != nil {
		return err
	} else {
		fmt.Println(state)
	}

	return nil
}
