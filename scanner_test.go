package libudev

import (
	"github.com/citilinkru/libudev/matcher"
	"testing"
)

func TestNewScanner(t *testing.T) {
	s := NewScanner()
	_, ok := interface{}(s).(*Scanner)
	if !ok {
		t.Fatal("Structure does not equal Scanner")
	}
}

func TestScanDevices(t *testing.T) {
	s := NewScanner()
	s.devicesPath = "./assets/fixtures/demo_tree/sys/devices"
	s.udevDataPath = "./assets/fixtures/demo_tree/run/udev/data"
	err, devices := s.ScanDevices()
	if err != nil {
		t.Fatal("Error scan demo tree")
	}

	if len(devices) != 11 {
		t.Fatal("Scanned devices count not equal 11")
	}

	m := matcher.NewMatcher()
	m.AddRule(matcher.NewRuleAttr("dev", "189:133"))
	dFiltered := m.Match(devices)
	if len(dFiltered) != 1 {
		t.Fatal("Not found device by Attr `dev` = `189:133`")
	}

	if len(dFiltered[0].Children) != 1 {
		t.Fatal("Device (`dev` = `189:133`) children count not equal 1")
	}

	if dFiltered[0].Children[0].Env["DEVNAME"] != "usb/lp0" {
		t.Fail()
	}

	if dFiltered[0].Parent == nil {
		t.Fatal("Not found parent device for device (`dev` = `189:133`)")
	}

	if dFiltered[0].Parent.Attrs["idProduct"] != "0024" {
		t.Fail()
	}
}

func TestScanDevicesIfNotSupported(t *testing.T) {
	s := NewScanner()
	s.devicesPath = "./NOT_EXIT_DIR"
	s.udevDataPath = "./NOT_EXIT_DIR"
	_, devices := s.ScanDevices()

	if len(devices) != 0 {
		t.Fatal("If the scan fails, then the device can not be found")
	}
}
