package matcher

import (
	"github.com/citilinkru/libudev/types"
	"regexp"
)

// RuleAttr structure of the filtering rule by attributes.
type RuleAttr struct {
	attrName string
	regexp   *regexp.Regexp
}

// NewRuleAttr creates a new instance of the filtering rule by attributes.
func NewRuleAttr(attrName, regexpValue string) *RuleAttr {
	rule := &RuleAttr{
		attrName: attrName,
	}

	r, err := regexp.Compile(regexpValue)
	if err == nil {
		rule.regexp = r
	}

	return rule
}

// Match verifies that the device complies with the rule.
//
// device - device for checking rules
func (m *RuleAttr) Match(device *types.Device) bool {
	if m.regexp == nil {
		return false
	}

	attrValue, ok := device.Attrs[m.attrName]
	if !ok {
		return false
	}

	return m.regexp.MatchString(attrValue)
}
