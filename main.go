package main

import (
	"github.com/jscherff/gocmdb/usbci/magtek"
	"github.com/jscherff/gocmdb"
	"github.com/google/gousb"
	//"encoding/json"
	"encoding/xml"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {

// =========================== TESTS ===========================

/*
	dj := "{\"HostName\":\"John-SurfacePro\",\"DeviceSN\":\"B164F78\",\"VendorID\":\"0801\",\"ProductID\":\"0001\",\"SoftwareID\":\"21042840G01\",\"VendorName\":\"Mag-Tek\",\"ProductName\":\"USB Swipe Reader\",\"ProductVer\":\"V05\",\"FactorySN\":\"B164F78022713AA\",\"DescriptSN\":\"B164F78\",\"BusNumber\":\"1\",\"BusAddress\":\"4\",\"USBSpec\":\"1.10\",\"USBClass\":\"per-interface\",\"USBSubclass\":\"per-interface\",\"USBProtocol\":\"0\",\"DeviceSpeed\":\"full\",\"DeviceVer\":\"1.00\",\"MaxPktSize\":\"8\",\"BufferSize\":\"60\"}"
	di := new(magtek.DeviceInfo)
	_ = json.Unmarshal([]byte(dj), di)
	fmt.Println(di)
*/
	var e error

	dx1 := []byte("<DeviceInfo><HostName>John-SurfacePro</HostName><DeviceSN>24FA12C</DeviceSN><VendorID>0801</VendorID><ProductID>0001</ProductID><SoftwareID>21042818B01</SoftwareID><ProductVer></ProductVer><FactorySN></FactorySN><VendorName>Mag-Tek</VendorName><ProductName>USB Swipe Reader</ProductName><DescriptSN>24FA12C</DescriptSN><BusNumber>1</BusNumber><BusAddress>8</BusAddress><USBSpec>1.10</USBSpec><USBClass>per-interface</USBClass><USBSubclass>per-interface</USBSubclass><USBProtocol>0</USBProtocol><DeviceSpeed>full</DeviceSpeed><DeviceVer>1.00</DeviceVer><MaxPktSize>8</MaxPktSize><BufferSize>24</BufferSize></DeviceInfo>")
	dx2 := []byte("<DeviceInfo><HostName>John-SurfacePro</HostName><DeviceSN>24FA12D</DeviceSN><VendorID>0801</VendorID><ProductID>0001</ProductID><SoftwareID>21042818B01</SoftwareID><ProductVer></ProductVer><FactorySN></FactorySN><VendorName>Mag-Tek</VendorName><ProductName>USB Swipe Reader</ProductName><DescriptSN>24FA12C</DescriptSN><BusNumber>1</BusNumber><BusAddress>8</BusAddress><USBSpec>1.10</USBSpec><USBClass>per-interface</USBClass><USBSubclass>per-interface</USBSubclass><USBProtocol>0</USBProtocol><DeviceSpeed>full</DeviceSpeed><DeviceVer>1.00</DeviceVer><MaxPktSize>8</MaxPktSize><BufferSize>24</BufferSize></DeviceInfo>")
	dx3 := []byte("<DeviceInfo><HostName>John-SurfacePro</HostName><DeviceSN>24FA12D</DeviceSN><VendorID>0801</VendorID><ProductID>0001</ProductID><SoftwareID>21042818B01</SoftwareID><FactorySN></FactorySN><VendorName>Mag-Tek</VendorName><ProductName>USB Swipe Reader</ProductName><DescriptSN>24FA12C</DescriptSN><BusNumber>1</BusNumber><BusAddress>8</BusAddress><USBSpec>1.10</USBSpec><USBClass>per-interface</USBClass><USBSubclass>per-interface</USBSubclass><USBProtocol>0</USBProtocol><DeviceSpeed>full</DeviceSpeed><DeviceVer>1.00</DeviceVer><MaxPktSize>8</MaxPktSize><BufferSize>24</BufferSize></DeviceInfo>")

	di1 := new(magtek.DeviceInfo)
	di2 := new(magtek.DeviceInfo)
	di3 := new(magtek.DeviceInfo)

	e = xml.Unmarshal(dx1, di1)
	if e != nil {log.Fatalf("%v", e)}
	e = xml.Unmarshal(dx2, di2)
	if e != nil {log.Fatalf("%v", e)}
	e = xml.Unmarshal(dx3, di3)
	if e != nil {log.Fatalf("%v", e)}

	dc1, e := di1.CSV(false)
	fmt.Println(string(dc1))

	dc2, e := di2.CSV(false)
	fmt.Println(string(dc2))

	dc3, e := di3.CSV(false)
	fmt.Println(string(dc3))

	ss, e := gocmdb.StructCompare(*di1, *di2)
	if e != nil {fmt.Printf("ERROR: %v", e)}
	fmt.Println(ss)

	ss, e = gocmdb.StructCompare(*di2, *di3)
	if e != nil {fmt.Printf("ERROR: %v", e)}
	fmt.Println(ss)

	os.Exit(0)
// =========================== TESTS ===========================

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "You must specify a mode of operation.\n")
		fsMode.Usage()
		os.Exit(1)
	}

	fsMode.Parse(os.Args[1:2])

	var flagset *flag.FlagSet

	switch {

	case *fModeReport:
		flagset = fsReport

	case *fModeConfig:
		flagset = fsConfig

	case *fModeReset:
		flagset = fsReset
	}

	if flagset.Parse(os.Args[2:]); flagset.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "You must specify at least one option.\n")
		flagset.Usage()
		os.Exit(1)
	}

	context := gousb.NewContext()
	defer context.Close()

	// Open devices that report a Magtek vendor ID, 0x0801.
	// We omit error checking on OpenDevices() because this
	// function terminates with 'libusb: not found [code -5]'
	// on Windows systems.

	devices, _ := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return uint16(desc.Vendor) == magtek.MagtekVendorID
	})

	if len(devices) == 0 {
		log.Fatalf("No devices found")
	}

	for _, device := range devices {

		defer device.Close()
		device, err := magtek.NewDevice(device)

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

// =========================== TESTS ===========================
/*
		di1, _ := magtek.NewDeviceInfo(device)
		b, _ := di1.JSON(false)

		fmt.Println(string(b))

		di2 := new(magtek.DeviceInfo)
		_ = json.Unmarshal(b, di2)

		fmt.Println(device)
		fmt.Println(di2)

		fmt.Println(di1.Matches(di2))

		os.Exit(0)
*/
// =========================== TESTS ===========================

		switch {

		case *fModeReport:
			err = report(device)

		case *fModeConfig:
			err = config(device)

		case *fModeReset:
			err = reset(device)
		}
	}
}
