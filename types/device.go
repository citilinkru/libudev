// Package types contains data structures
package types

// Device structure describing the device.
type Device struct {
	Devpath         string
	Env             map[string]string
	Attrs           map[string]string
	Tags            []string
	UsecInitialized string

	Parent   *Device
	Children []*Device
}
