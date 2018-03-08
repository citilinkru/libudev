/*
Package libudev implements a native udev library.

Recursively passes `/sys/devices/...` directory and reads `uevent` files (placed in `Env`).
Based on the information received from `uevent` files, it tries to enrich the data based
on the data received from the files that are on the same level as the `uevent` file (they are placed in `Attrs`),
and also tries to find and read the files `/run/udev/data/...` (placed in `Env` or `Tags`).

After building a list of devices, the library builds a device tree.
*/
package libudev

import (
	"bufio"
	"bytes"
	"github.com/citilinkru/libudev/types"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Scanner structure of the device scanner.
type Scanner struct {
	devicesPath  string
	udevDataPath string
}

// NewScanner creates a new instance of device scanner.
func NewScanner() *Scanner {
	return &Scanner{
		devicesPath:  "/sys/devices",
		udevDataPath: "/run/udev/data",
	}
}

// ScanDevices scans directories for `uevent` files and creates a device tree.
func (s *Scanner) ScanDevices() (err error, devices []*types.Device) {
	devices = []*types.Device{}
	devicesMap := map[string]*types.Device{}
	err = filepath.Walk(s.devicesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() || info.Name() != "uevent" {
			return nil
		}

		device := &types.Device{
			Devpath: filepath.Dir(path),
			Env:     map[string]string{},
			Attrs:   map[string]string{},
			Parent:  nil,
		}

		s.readAttrs(filepath.Dir(path), device)

		err = s.readUeventFile(path, device)
		if err != nil {
			return nil
		}

		devicesMap[device.Devpath] = device
		return nil
	})

	// make tree
	for _, v := range devicesMap {
		parts := strings.Split(v.Devpath, "/")

		devpath := v.Devpath
		for i := len(parts) - 1; i >= 0; i-- {
			devpath = strings.TrimSuffix(devpath, "/"+parts[i])
			if devpath == s.devicesPath {
				break
			}

			if device, ok := devicesMap[devpath]; ok {
				v.Parent = device
				device.Children = append(device.Children, v)
				break
			}
		}

		devices = append(devices, v)
	}

	return err, devices
}

func (s *Scanner) readAttrs(path string, device *types.Device) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() || f.Name() == "uevent" || f.Name() == "descriptors" {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			continue
		}

		device.Attrs[f.Name()] = strings.Trim(string(data), "\n\r\t ")
	}

	return nil
}

func (s *Scanner) readUeventFile(path string, device *types.Device) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	buf := bufio.NewScanner(bytes.NewBuffer(data))

	var line string
	for buf.Scan() {
		line = buf.Text()
		field := strings.SplitN(line, "=", 2)
		if len(field) != 2 {
			continue
		}

		device.Env[field[0]] = field[1]

	}

	devPath := filepath.Join(filepath.Dir(path), "dev")
	devString, err := s.readDevFile(devPath)
	if err != nil {
		return err
	}

	err = s.readUdevInfo(devString, device)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scanner) readDevFile(path string) (data string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}

	defer f.Close()

	d, err := ioutil.ReadAll(f)
	return strings.Trim(string(d), "\n\r\t "), err
}

func (s *Scanner) readUdevInfo(devString string, d *types.Device) error {
	path := fmt.Sprintf("%s/c%s", s.udevDataPath, devString)
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	buf := bufio.NewScanner(bytes.NewBuffer(data))

	var line string
	for buf.Scan() {
		line = buf.Text()
		groups := strings.SplitN(line, ":", 2)
		if len(groups) != 2 {
			continue
		}

		if groups[0] == "I" {
			d.UsecInitialized = groups[1]
			continue
		}

		if groups[0] == "G" {
			d.Tags = append(d.Tags, groups[1])
			continue
		}

		if groups[0] == "E" {
			fields := strings.SplitN(groups[1], "=", 2)
			if len(fields) != 2 {
				continue
			}

			d.Env[fields[0]] = fields[1]
		}
	}

	return nil
}
