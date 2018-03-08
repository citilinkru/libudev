package matcher

import (
	"github.com/citilinkru/libudev/types"
	"testing"
)

func TestNewRuleAttr(t *testing.T) {
	r := NewRuleAttr("TEST", "TEST")
	_, ok := interface{}(r).(Rule)
	if !ok {
		t.Fatal("Structure does not implement interface")
	}
}

func TestMatchAttr(t *testing.T) {
	dv1 := &types.Device{
		Attrs: map[string]string{"ATTR_1": "abc123def", "ATTR_2": "test data"},
	}

	r1 := NewRuleAttr("ATTR_1", "^abc[0-9]+")
	if !r1.Match(dv1) {
		t.Fatal("Could not find device `dv1`")
	}

	r2 := NewRuleAttr("ATTR_2", "^abc[0-9]+")
	if r2.Match(dv1) {
		t.Fatal("The device `dv1` was found incorrectly")
	}

	r3 := NewRuleAttr("ATTR_2", "test data")
	if !r3.Match(dv1) {
		t.Fatal("Could not find device `dv1`")
	}

	r4 := NewRuleAttr("ATTR_NOT_EXIST", "test data")
	if r4.Match(dv1) {
		t.Fatal("The device `dv1` was found incorrectly")
	}

	r5 := RuleAttr{regexp: nil, attrName: "ATTR_1"}
	if r5.Match(dv1) {
		t.Fatal("The device `dv1` was found incorrectly")
	}
}
