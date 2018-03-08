# libudev
Golang native implementation Udev library

[![Build Status](https://travis-ci.org/citilinkru/libudev.svg?branch=master)](https://travis-ci.org/citilinkru/libudev)
[![Coverage Status](https://coveralls.io/repos/github/citilinkru/libudev/badge.svg?branch=master)](https://coveralls.io/github/citilinkru/libudev?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/citilinkru/libudev)](https://goreportcard.com/report/github.com/citilinkru/libudev)
[![GoDoc](https://godoc.org/github.com/citilinkru/libudev?status.svg)](https://godoc.org/github.com/citilinkru/libudev)
[![GitHub release](https://img.shields.io/github/release/citilinkru/libudev.svg)](https://github.com/citilinkru/libudev/releases)


Installation
------------
    go get github.com/citilinkru/libudev

Usage
-----

### Scanning devices
```go
sc := libudev.NewScanner()
err, devices := s.ScanDevices()
```

### Filtering devices
```go
m := matcher.NewMatcher()
m.SetStrategy(matcher.StrategyOr)
m.AddRule(matcher.NewRuleAttr("dev", "189:133"))
m.AddRule(matcher.NewRuleEnv("DEVNAME", "usb/lp0"))

filteredDevices := m.Match(devices)
```

### Getting parent device
```go
if device.Parent != nil {
    fmt.Printf("%s\n", device.Parent.Devpath)
}
```

### Getting children devices
```go
fmt.Printf("Count children devices %d\n", len(device.Children))
```

Features
--------
* 100% Native code
* Without external dependencies
* Code is covered by tests

Requirements
------------

* Need at least `go1.8` or newer.

Documentation
-------------

You can read package documentation [here](http:godoc.org/github.com/citilinkru/libudev) or read tests.

Testing
-------
Unit-tests:
```bash
go test -race -v ./...
```

Contributing
------------
* Fork
* Write code
* Run unit test: `go test -race -v ./...`
* Run go vet: `go vet -v ./...`
* Run go fmt: `go fmt ./...`
* Commit changes
* Create pull-request
