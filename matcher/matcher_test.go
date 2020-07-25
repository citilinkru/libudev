package matcher

import (
	"github.com/citilinkru/libudev/types"
	"testing"
)

func TestNewMatcher(t *testing.T) {
	m := NewMatcher()
	_, ok := interface{}(m).(*Matcher)
	if !ok {
		t.Fatal("Structure does not equal Matcher")
	}
}

func TestSetStrategy(t *testing.T) {
	m := NewMatcher()
	m.SetStrategy(StrategyAnd)

	if m.strategy != StrategyAnd {
		t.Fatal("Failed to set strategy:", StrategyAnd)
	}

	m.SetStrategy(StrategyOr)

	if m.strategy != StrategyOr {
		t.Fatal("Failed to set strategy:", StrategyOr)
	}
}

func TestAddRule(t *testing.T) {
	m := NewMatcher()
	m.AddRule(NewRuleDevpath("TEST1"))
	m.AddRule(NewRuleDevpath("TEST2"))
	if len(m.rules) != 2 {
		t.Fatal("Failed to add two rules")
	}
}

func TestMatchDefault(t *testing.T) {
	devices := getDemoDevices()

	m := NewMatcher()
	m.AddRule(NewRuleDevpath("devpaht-1"))
	m.AddRule(NewRuleDevpath("devpaht-2"))
	if len(m.Match(devices)) != 0 {
		t.Fatal("Devices was found incorrectly")
	}

	m2 := NewMatcher()
	m2.AddRule(NewRuleDevpath("devpaht-.+"))
	m2.AddRule(NewRuleDevpath("devpaht-[0-9]+"))
	if len(m2.Match(devices)) != 2 {
		t.Fatal("Not found two devices")
	}

	m3 := NewMatcher()
	if len(m3.Match(devices)) != 0 {
		t.Fatal("Empty rules Matcher finded not 0 devices")
	}
}

func TestMatchAnd(t *testing.T) {
	devices := getDemoDevices()

	m := NewMatcher()
	m.SetStrategy(StrategyAnd)
	m.AddRule(NewRuleDevpath("devpaht-1"))
	m.AddRule(NewRuleDevpath("devpaht-2"))
	if len(m.Match(devices)) != 0 {
		t.Fatal("Devices was found incorrectly")
	}

	m2 := NewMatcher()
	m2.SetStrategy(StrategyAnd)
	m2.AddRule(NewRuleDevpath("devpaht-.+"))
	m2.AddRule(NewRuleDevpath("devpaht-[0-9]+"))
	if len(m2.Match(devices)) != 2 {
		t.Fatal("Not found two devices")
	}

	m2.AddRule(NewRuleEnv("ENV-1", "[a-z]+"))
	m2.AddRule(NewRuleEnv("ENV-2", "[0-9]+"))
	m2.AddRule(NewRuleAttr("ATTR-1", "[a-z]+"))
	m2.AddRule(NewRuleAttr("ATTR-2", "[0-9]+"))
	if len(m2.Match(devices)) != 2 {
		t.Fatal("Not found two devices")
	}

	m2.AddRule(NewRuleAttr("ATTR-2", "123"))
	if len(m.Match(devices)) != 0 {
		t.Fatal("Devices was found incorrectly")
	}

	m3 := NewMatcher()
	m3.SetStrategy(StrategyAnd)
	if len(m3.Match(devices)) != 0 {
		t.Fatal("Empty rules Matcher finded not 0 devices")
	}
}

func TestMatchOr(t *testing.T) {
	devices := getDemoDevices()

	m := NewMatcher()
	m.SetStrategy(StrategyOr)
	m.AddRule(NewRuleDevpath("devpaht-1"))
	m.AddRule(NewRuleDevpath("devpaht-3"))
	if len(m.Match(devices)) != 1 {
		t.Fatal("Not found one device")
	}

	m.AddRule(NewRuleDevpath("devpaht-2"))
	if len(m.Match(devices)) != 2 {
		t.Fatal("Not found two devices")
	}

	m2 := NewMatcher()
	m2.SetStrategy(StrategyOr)
	m2.AddRule(NewRuleDevpath("devpaht-1"))
	m2.AddRule(NewRuleEnv("ENV-1", "ghi456jkl"))
	m2.AddRule(NewRuleEnv("ENV-2", "NOT_FOUND"))
	m2.AddRule(NewRuleAttr("ATTR-1", "ghi456jkl"))
	m2.AddRule(NewRuleAttr("ATTR-2", "NOT_FOUND"))
	if len(m2.Match(devices)) != 2 {
		t.Fatal("Not found two devices")
	}

	m3 := NewMatcher()
	m3.SetStrategy(StrategyOr)
	if len(m3.Match(devices)) != 0 {
		t.Fatal("Empty rules Matcher finded not 0 devices")
	}
}

func getDemoDevices() []*types.Device {
	return []*types.Device{
		{
			Devpath: "devpaht-1",
			Env:     map[string]string{"ENV-1": "abc123def", "ENV-2": "123"},
			Attrs:   map[string]string{"ATTR-1": "abc123def", "ATTR-2": "123"},
		},
		{
			Devpath: "devpaht-2",
			Env:     map[string]string{"ENV-1": "ghi456jkl", "ENV-2": "456"},
			Attrs:   map[string]string{"ATTR-1": "ghi456jkl", "ATTR-2": "456"},
		},
	}
}
