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
	"strings"
	"flag"
	"fmt"

	"github.com/jscherff/gocmdb"
)

const (
	nameIndex int = 0
	valueIndex int = 1
)

var (
	fsAction = flag.NewFlagSet("action", flag.ExitOnError)
	fActionAudit = fsAction.Bool("audit", false, "Audit devices")
	fActionCheckin = fsAction.Bool("checkin", false, "Check devices in")
	fActionLegacy = fsAction.Bool("legacy", false, "Legacy operation")
	fActionReport = fsAction.Bool("report", false, "Report actions")
	fActionReset = fsAction.Bool("reset", false, "Reset device")
	fActionSerial = fsAction.Bool("serial", false, "Set serial number")

	fsReport = flag.NewFlagSet("report", flag.ExitOnError)
	fReportFile = fsReport.String("file", "", "Write report to `<file>`")
	fReportStdout = fsReport.Bool("stdout", false, "Write output to stdout")
	fReportFormat *string

	fsSerial = flag.NewFlagSet("serial", flag.ExitOnError)
	fSerialCopy = fsSerial.Bool("copy", false, "Copy factory serial number")
	fSerialErase = fsSerial.Bool("erase", false, "Erase current serial number")
	fSerialForce = fsSerial.Bool("force", false, "Force serial number change")
	fSerialServer = fsSerial.Bool("server", false, "Set serial number from server")
	fSerialConfig = fsSerial.String("config", "", "Set serial number to `<string>`")
)

func init() {
	var formats []string
	for _, f := range gocmdb.ReportFormats {formats = append(formats, f[nameIndex])}
	usage := fmt.Sprintf("`<fmt>` = {%s}", strings.Join(formats, "|"))
	fReportFormat = fsReport.String("format", "csv", usage)
}
