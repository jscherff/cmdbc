package main

import (
	"github.com/jscherff/gocmdb"
	"strings"
	"flag"
	"fmt"
)

var (
	fsAction = flag.NewFlagSet("action", flag.ExitOnError)
	fActionReport = fsAction.Bool("report", false, "Generate Report")
	fActionConfig = fsAction.Bool("config", false, "Configure device")
	fActionReset = fsAction.Bool("reset", false, "Reset device")
	fActionAudit = fsAction.Bool("audit", false, "Audit device")
	fActionCheckin = fsAction.String("checkin", "", "Checkin to `<url>`")

	fsReport = flag.NewFlagSet("report", flag.ExitOnError)
	fReportFile = fsReport.String("file", "", "Write report to `<file>`")
	fReportStdout = fsReport.Bool("stdout", false, "Write output to stdout")
	fReportFormat *string

	fsConfig = flag.NewFlagSet("config", flag.ExitOnError)
	fConfigCopy = fsConfig.Bool("copy", false, "Copy factory serial number")
	fConfigErase = fsConfig.Bool("erase", false, "Erase current serial number")
	fConfigForce = fsConfig.Bool("force", false, "Force serial number change")
	fConfigString = fsConfig.String("string", "", "Set serial number to `<value>`")
	fConfigServer = fsConfig.String("server", "", "Set serial number from `<url>`")
)

func init() {
	var formats []string
	for _, f := range gocmdb.ReportFormats {formats = append(formats, f[NameIx])}
	usage := fmt.Sprintf("Where `<fmt>` is {%s}", strings.Join(formats, "|"))
	fReportFormat = fsReport.String("format", "csv", usage)
}
