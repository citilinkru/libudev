package libudev

import (
	"archive/zip"
	"github.com/citilinkru/libudev/matcher"
	"io"
	"os"
	"path/filepath"
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
	Unzip("./assets/fixtures/demo_tree.zip", "/tmp/demo_tree/")
	s.devicesPath = "/tmp/demo_tree/demo_tree/sys/devices"
	s.udevDataPath = "/tmp/demo_tree/demo_tree/run/udev/data"
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

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
