package main

import (
	"github.com/jscherff/gocmdb"
	"flag"
	"fmt"
)

var (
	fsMode = flag.NewFlagSet("mode", flag.ExitOnError)
	fModeReport = fsMode.Bool("report", false, "Report mode")
	fModeConfig = fsMode.Bool("config", false, "Config mode")
	fModeReset = fsMode.Bool("reset", false, "Reset mode")
)

var (
	fsReport = flag.NewFlagSet("report", flag.ExitOnError)
	fReportAll = fsReport.Bool("all", false, "Include all report fields")
	fReportFile = fsReport.String("file", "", "Write report to `<file>`")
	fReportServer = fsReport.String("server", "", "Submit report to server `<url>`")
	fReportStdout = fsReport.Bool("stdout", false, "Write output to stdout")
	fReportFormat *string
)

var (
	fsConfig = flag.NewFlagSet("config", flag.ExitOnError)
	fConfigCopy = fsConfig.Bool("copy", false, "Copy factory serial number")
	fConfigErase = fsConfig.Bool("erase", false, "Erase current serial number")
	fConfigForce = fsConfig.Bool("force", false, "Force serial number change")
	fConfigString = fsConfig.String("string", "", "Set serial number to string `<value>`")
	fConfigServer = fsConfig.String("server", "", "Set serial number from server `<url>`")

)

var (
	fsReset = flag.NewFlagSet("reset", flag.ExitOnError)
	fResetUsb = fsReset.Bool("usb", false, "Perform a USB reset")
	fResetDev = fsReset.Bool("dev", false, "Perform a device reset")
)

func init() {

	gocmdb.ReportFormats = append(gocmdb.ReportFormats, []string{"leg", "Legacy report format"})

	usage := "Where `<format>` is one of:"
	for _, f := range gocmdb.ReportFormats {
		usage += fmt.Sprintf("\n\t%q\t%s", f[gocmdb.NameIx], f[gocmdb.ValueIx])
	}

	fReportFormat = fsReport.String("format", "", usage)
}
