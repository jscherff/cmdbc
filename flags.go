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
	fReportFile = fsReport.String ("file", "", "Write report to `<file>`")
	fReportStdout = fsReport.Bool("stdout", false, "Write output to stdout")
	fReportFormat *string
)

var (
	fsConfig = flag.NewFlagSet("config", flag.ExitOnError)
	fConfigCopy = fsConfig.Bool("copy", false, "Copy factory serial number")
	fConfigErase = fsConfig.Bool("erase", false, "Erase current serial number")
	fConfigForce = fsConfig.Bool("force", false, "Force serial number change")
	fConfigSet = fsConfig.String("set", "", "Set serial number to `<string>`")
	fConfigUrl = fsConfig.String("url", "", "Set serial number from `<url>`")

)

var (
	fsReset = flag.NewFlagSet("reset", flag.ExitOnError)
	fResetUsb = fsReset.Bool("usb", false, "Perform a USB reset")
	fResetDev = fsReset.Bool("dev", false, "Perform a device reset")
)

func init() {
	usage := "Report format `<format>`"
	for k, v := range gocmdb.ReportFormats {usage += fmt.Sprintf("\n\t%q\t%s", k, v)}
	fReportFormat = fsReport.String("format", "csv", usage)
}
