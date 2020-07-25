/*
Package matcher implements a service for filtering devices.

Matcher allows you to add multiple filter rules by using the method `AddRule`
and effectively filter the list of devices in a single pass.

By default, Matcher uses the `AND` comparison strategy, but you can set the filtering strategy to` OR`.

*/
package matcher

import (
	"github.com/citilinkru/libudev/types"
)

const (
	StrategyOr  = "OR"
	StrategyAnd = "AND"
)

// Rule interface rules for filtering devices.
type Rule interface {
	Match(device *types.Device) bool
}

// Matcher structure of the device filter.
type Matcher struct {
	rules    []Rule
	strategy string
}

// NewMatcher creates a new Matcher instance.
func NewMatcher() *Matcher {
	return &Matcher{
		rules:    []Rule{},
		strategy: StrategyAnd,
	}
}

// SetStrategy sets the device filtering strategy.
//
// strategy - filtering strategy, see the constants `StrategyOr` and` StrategyAnd`
func (m *Matcher) SetStrategy(strategy string) {
	m.strategy = strategy
}

// AddRule adds a new device filtering rule.
//
// rule - device filtering rule
func (m *Matcher) AddRule(rule Rule) {
	m.rules = append(m.rules, rule)
}

func (m *Matcher) Match(devices []*types.Device) []*types.Device {
	var ret []*types.Device
	for _, v := range devices {
		if !m.matchDevice(v) {
			continue
		}

		ret = append(ret, v)
	}

	return ret
}

func (m *Matcher) matchDevice(device *types.Device) bool {
	if len(m.rules) == 0 {
		return false
	}

	def := false
	if m.strategy == StrategyAnd {
		def = true
	}

	for _, v := range m.rules {
		if m.strategy == StrategyAnd && !v.Match(device) {
			return false
		}

		if m.strategy == StrategyOr && v.Match(device) {
			return true
		}
	}

	return def
}
