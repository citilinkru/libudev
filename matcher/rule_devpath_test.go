package matcher

import (
	"github.com/citilinkru/libudev/types"
	"testing"
)

func TestNewRuleDevpath(t *testing.T) {
	r := NewRuleDevpath("TEST")
	_, ok := interface{}(r).(Rule)
	if !ok {
		t.Fatal("Structure does not implement interface")
	}
}

func TestMatchDevpath(t *testing.T) {
	dv1 := &types.Device{
		Devpath: "/sys/devices/test1",
	}

	r1 := NewRuleDevpath("test[0-9]+$")
	if !r1.Match(dv1) {
		t.Fatal("Could not find device `dv1`")
	}

	r2 := NewRuleDevpath("^test[0-9]+$")
	if r2.Match(dv1) {
		t.Fatal("The device `dv1` was found incorrectly")
	}

	r3 := NewRuleDevpath("/sys/devices/test1")
	if !r3.Match(dv1) {
		t.Fatal("Could not find device `dv1`")
	}

	r4 := RuleDevpath{regexp: nil}
	if r4.Match(dv1) {
		t.Fatal("The device `dv1` was found incorrectly")
	}
}
