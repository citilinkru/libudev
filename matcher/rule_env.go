package matcher

import (
	"github.com/citilinkru/libudev/types"
	"regexp"
)

// RuleEnv structure of the filtering rule by `Env`.
type RuleEnv struct {
	envName string
	regexp  *regexp.Regexp
}

// NewRuleEnv creates a new instance of the filtering rule by `Env`.
func NewRuleEnv(envName, regexpValue string) *RuleEnv {
	rule := &RuleEnv{
		envName: envName,
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
func (m *RuleEnv) Match(device *types.Device) bool {
	if m.regexp == nil {
		return false
	}

	envValue, ok := device.Env[m.envName]
	if !ok {
		return false
	}

	return m.regexp.MatchString(envValue)
}
