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
	//`log`
	//`os`
	`io/ioutil`
	`fmt`
	`path/filepath`
	`crypto/sha256`
	`testing`
	`github.com/jscherff/gocmdb/usbci`
)

func TestConfig(t *testing.T) {

	var err error

	if conf, err = NewConfig(`config.json`); err == nil {
		slog, clog, elog = NewLoggers()
	} else {
		t.Error(err)
	}
}

func TestRestore(t *testing.T) {

	var err error

	t.Run("Magek Restore 1", func(t *testing.T) {
		if mag1, err = usbci.NewMagtek(nil); err == nil {
			err = mag1.RestoreJSON(magJSNv1)
		}
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Magek Restore 2", func(t *testing.T) {
		if mag2, err = usbci.NewMagtek(nil); err == nil {
			err = mag2.RestoreJSON(magJSNv2)
		}
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Generic Restore 1", func(t *testing.T) {
		if gen1, err = usbci.NewGeneric(nil); err == nil {
			err = gen1.RestoreJSON(genJSNv1)
		}
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Generic Restore 2", func(t *testing.T) {
		if gen2, err = usbci.NewGeneric(nil); err == nil {
			err = gen2.RestoreJSON(genJSNv2)
		}
		if err != nil {
			t.Error(err)
		}
	})
}

func TestReport(t *testing.T) {

	var (
		b []byte
		fn string
		err error
	)

	*fReportConsole = false

	t.Run("JSN Report", func(t *testing.T) {
		fn = filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.json`)
		*fReportFormat = `json`

		if err = reportAction(mag1); err == nil {
			if b, err = ioutil.ReadFile(fn); err == nil {
				if mag1ShaJSN != sha256.Sum256(b) {
					t.Error(`Sha256 hash of JSON report is incorrect`)
				}
			}
		}
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("XML Report", func(t *testing.T) {
		fn = filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.xml`)
		*fReportFormat = `xml`

		if err = reportAction(mag1); err == nil {
			if b, err = ioutil.ReadFile(fn); err == nil {
				if mag1ShaXML != sha256.Sum256(b) {
					t.Error(`Sha256 hash of XML report is incorrect`)
				}
			}
		}
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("CSV Report", func(t *testing.T) {
		fn = filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.csv`)
		*fReportFormat = `csv`

		if err = reportAction(mag1); err == nil {
			if b, err = ioutil.ReadFile(fn); err == nil {
				if mag1ShaCSV != sha256.Sum256(b) {
					t.Error(`Sha256 hash of CSV report is incorrect`)
				}
			}
		}
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("NVP Report", func(t *testing.T) {
		fn = filepath.Join(conf.Paths.ReportDir, mag1.Filename() + `.nvp`)
		*fReportFormat = `nvp`

		if err = reportAction(mag1); err == nil {
			if b, err = ioutil.ReadFile(fn); err == nil {
				if mag1ShaNVP != sha256.Sum256(b) {
					t.Error(`Sha256 hash of NVP report is incorrect`)
				}
			}
		}
		if err != nil {
			t.Error(err)
		}
	})
}

func TestCheckinDevice(t *testing.T) {

	var (
		j []byte
		err error
		ss [][]string
		mag3 *usbci.Magtek
	)


	if mag3, err = usbci.NewMagtek(nil); err != nil {
		t.Error(err)
		return
	}

	t.Run("Checkin and Checkout Match", func(t *testing.T) {
		if err = CheckinDevice(mag1); err == nil {
			if j, err = CheckoutDevice(mag1); err == nil {
				if ss, err = mag1.CompareJSON(j); err == nil {
					if len(ss) != 0 {
						err = fmt.Errorf(`checkin and checkout mismatch`)
					}
				}
			}
		}
		if err != nil {
			t.Error(err)
		}
	})

	//if err = CheckinDevice(mag2); err != nil {
	//	t.Error(err)
	//}
}

//func TestGetNewSN(t *testing.T) {

/*
func legacyAction(o gocmdb.Reportable) (err error) {
func reportAction(o gocmdb.Reportable) (err error) {
func serialAction(o gocmdb.Configurable) (err error) {
func auditAction(o gocmdb.Auditable) (err error) {
func GetNewSN(o gocmdb.Registerable) (s string, err error) {
func CheckinDevice(o gocmdb.Registerable) (err error) {
func CheckoutDevice(o gocmdb.Auditable) (j []byte, err error) {
func SubmitAudit(o gocmdb.Auditable) (err error) {
func httpPost(url string, j []byte ) (b []byte, sc int, err error) {
func httpGet(url string) (b []byte, sc int, err error) {
func httpRequest(req *http.Request) (b []byte, sc int, err error) {

/*
	fsAction = flag.NewFlagSet("action", flag.ExitOnError)
	fActionAudit = fsAction.Bool("audit", false, "Audit devices")
	fActionCheckin = fsAction.Bool("checkin", false, "Check devices in")
	fActionLegacy = fsAction.Bool("legacy", false, "Legacy operation")
	fActionReport = fsAction.Bool("report", false, "Report actions")
	fActionReset = fsAction.Bool("reset", false, "Reset device")
	fActionSerial = fsAction.Bool("serial", false, "Set serial number")

	fsReport = flag.NewFlagSet("report", flag.ExitOnError)
	fReportFolder = fsReport.String("folder", "", "Write reports to `<path>`")
	fReportConsole = fsReport.Bool("console", false, "Write reports to console")
	fReportFormat = fsReport.String("format", "csv", "Report `<format>` {csv|nvp|xml|json}")

	fsSerial = flag.NewFlagSet("serial", flag.ExitOnError)
	fSerialCopy = fsSerial.Bool("copy", false, "Copy factory serial number")
	fSerialErase = fsSerial.Bool("erase", false, "Erase current serial number")
	fSerialForce = fsSerial.Bool("force", false, "Force serial number change")
	fSerialFetch = fsSerial.Bool("fetch", false, "Fetch serial number from server")
	fSerialSet = fsSerial.String("set", "", "Set serial number to `<string>`")

	fsAudit = flag.NewFlagSet("audit", flag.ExitOnError)
	fAuditLocal = fsAudit.Bool("local", false, "Audit against local state")
	fAuditServer = fsAudit.Bool("server", false, "Audit against server state")
*/

/*
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

package gocmdb

type Identifiable interface {
	ID() (string)
	VID() (string)
	PID() (string)
	Host() (string)
	Type() (string)
	JSN() ([]byte, error)
	XML() ([]byte, error)
	CSV() ([]byte, error)
	NVP() ([]byte, error)
	Legacy() ([]byte)
	PrettyJSN() ([]byte, error)
	PrettyXML() ([]byte, error)
	Filename() (string)
	Save(string) (error)
	RestoreFile(string) (error)
	RestoreJSON([]byte) (error)
	CompareFile(string) ([][]string, error)
	CompareJSON([]byte) ([][]string, error)
	AuditFile(string) (error)
	AuditJSON([]byte) (error)
	AddChange(string, string, string)
	SetChanges([][]string)
	GetChanges() ([][]string)
	Matches(interface{}) (bool)
	Refresh() (map[string]bool)
	Reset() (error)
	CheckoutDeviceSN() (string, error)
	SetDeviceSN(string) (error)
	EraseDeviceSN() (error)
	SetFactorySN(string) (error)
	CopyFactorySN(int) (error)
}


type tester struct {
	flag    int
	prefix  string
	pattern string // regexp that log output must match; we add ^ and expected_text$ always
}

var tests = []tester{
	// individual pieces:
	{0, "", ""},
	{0, "XXX", "XXX"},
	{Ldate, "", Rdate + " "},
	{Ltime, "", Rtime + " "},
	{Ltime | Lmicroseconds, "", Rtime + Rmicroseconds + " "},
	{Lmicroseconds, "", Rtime + Rmicroseconds + " "}, // microsec implies time
	{Llongfile, "", Rlongfile + " "},
	{Lshortfile, "", Rshortfile + " "},
	{Llongfile | Lshortfile, "", Rshortfile + " "}, // shortfile overrides longfile
	// everything at once:
	{Ldate | Ltime | Lmicroseconds | Llongfile, "XXX", "XXX" + Rdate + " " + Rtime + Rmicroseconds + " " + Rlongfile + " "},
	{Ldate | Ltime | Lmicroseconds | Lshortfile, "XXX", "XXX" + Rdate + " " + Rtime + Rmicroseconds + " " + Rshortfile + " "},
}

// Test using Println("hello", 23, "world") or using Printf("hello %d world", 23)
func testPrint(t *testing.T, flag int, prefix string, pattern string, useFormat bool) {
	buf := new(bytes.Buffer)
	SetOutput(buf)
	SetFlags(flag)
	SetPrefix(prefix)
	if useFormat {
		Printf("hello %d world", 23)
	} else {
		Println("hello", 23, "world")
	}
	line := buf.String()
	line = line[0 : len(line)-1]
	pattern = "^" + pattern + "hello 23 world$"
	matched, err4 := regexp.MatchString(pattern, line)
	if err4 != nil {
		t.Fatal("pattern did not compile:", err4)
	}
	if !matched {
		t.Errorf("log output should match %q is %q", pattern, line)
	}
	SetOutput(os.Stderr)
}

func TestAll(t *testing.T) {
	for _, testcase := range tests {
		testPrint(t, testcase.flag, testcase.prefix, testcase.pattern, false)
		testPrint(t, testcase.flag, testcase.prefix, testcase.pattern, true)
	}
}

func TestOutput(t *testing.T) {
	const testString = "test"
	var b bytes.Buffer
	l := New(&b, "", 0)
	l.Println(testString)
	if expect := testString + "\n"; b.String() != expect {
		t.Errorf("log output should match %q is %q", expect, b.String())
	}
}

func TestFlagAndPrefixSetting(t *testing.T) {
	var b bytes.Buffer
	l := New(&b, "Test:", LstdFlags)
	f := l.Flags()
	if f != LstdFlags {
		t.Errorf("Flags 1: expected %x got %x", LstdFlags, f)
	}
	l.SetFlags(f | Lmicroseconds)
	f = l.Flags()
	if f != LstdFlags|Lmicroseconds {
		t.Errorf("Flags 2: expected %x got %x", LstdFlags|Lmicroseconds, f)
	}
	p := l.Prefix()
	if p != "Test:" {
		t.Errorf(`Prefix: expected "Test:" got %q`, p)
	}
	l.SetPrefix("Reality:")
	p = l.Prefix()
	if p != "Reality:" {
		t.Errorf(`Prefix: expected "Reality:" got %q`, p)
	}
	// Verify a log message looks right, with our prefix and microseconds present.
	l.Print("hello")
	pattern := "^Reality:" + Rdate + " " + Rtime + Rmicroseconds + " hello\n"
	matched, err := regexp.Match(pattern, b.Bytes())
	if err != nil {
		t.Fatalf("pattern %q did not compile: %s", pattern, err)
	}
	if !matched {
		t.Error("message did not match pattern")
	}
}

func TestUTCFlag(t *testing.T) {
	var b bytes.Buffer
	l := New(&b, "Test:", LstdFlags)
	l.SetFlags(Ldate | Ltime | LUTC)
	// Verify a log message looks right in the right time zone. Quantize to the second only.
	now := time.Now().UTC()
	l.Print("hello")
	want := fmt.Sprintf("Test:%d/%.2d/%.2d %.2d:%.2d:%.2d hello\n",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	got := b.String()
	if got == want {
		return
	}
	// It's possible we crossed a second boundary between getting now and logging,
	// so add a second and try again. This should very nearly always work.
	now = now.Add(time.Second)
	want = fmt.Sprintf("Test:%d/%.2d/%.2d %.2d:%.2d:%.2d hello\n",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	if got == want {
		return
	}
	t.Errorf("got %q; want %q", got, want)
}

func TestEmptyPrintCreatesLine(t *testing.T) {
	var b bytes.Buffer
	l := New(&b, "Header:", LstdFlags)
	l.Print()
	l.Println("non-empty")
	output := b.String()
	if n := strings.Count(output, "Header"); n != 2 {
		t.Errorf("expected 2 headers, got %d", n)
	}
	if n := strings.Count(output, "\n"); n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func BenchmarkItoa(b *testing.B) {
	dst := make([]byte, 0, 64)
	for i := 0; i < b.N; i++ {
		dst = dst[0:0]
		itoa(&dst, 2015, 4)   // year
		itoa(&dst, 1, 2)      // month
		itoa(&dst, 30, 2)     // day
		itoa(&dst, 12, 2)     // hour
		itoa(&dst, 56, 2)     // minute
		itoa(&dst, 0, 2)      // second
		itoa(&dst, 987654, 6) // microsecond
	}
}

func BenchmarkPrintln(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf, "", LstdFlags)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

func BenchmarkPrintlnNoFlags(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf, "", 0)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(testString)
	}
}

*/
