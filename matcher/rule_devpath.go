package matcher

import (
	"github.com/citilinkru/libudev/types"
	"regexp"
)

// RuleDevpath structure of the filtering rule by `Devpath`.
type RuleDevpath struct {
	regexp *regexp.Regexp
}

// NewRuleDevpath creates a new instance of the filtering rule by `Devpath`.
func NewRuleDevpath(regexpValue string) *RuleDevpath {
	rule := &RuleDevpath{}
	r, err := regexp.Compile(regexpValue)
	if err == nil {
		rule.regexp = r
	}

	return rule
}

// Match verifies that the device complies with the rule.
//
// device - device for checking rules
func (m *RuleDevpath) Match(device *types.Device) bool {
	if m.regexp == nil {
		return false
	}

	return m.regexp.MatchString(device.Devpath)
}
