package matcher

import (
	"github.com/citilinkru/libudev/types"
	"testing"
)

func TestNewRuleEnv(t *testing.T) {
	r := NewRuleEnv("TEST", "TEST")
	_, ok := interface{}(r).(Rule)
	if !ok {
		t.Fatal("Structure does not implement interface")
	}
}

func TestMatchEnv(t *testing.T) {
	dv1 := &types.Device{
		Env: map[string]string{"ENV_1": "abc123def", "ENV_2": "test data"},
	}

	r1 := NewRuleEnv("ENV_1", "^abc[0-9]+")
	if !r1.Match(dv1) {
		t.Fatal("Could not find device `dv1`")
	}

	r2 := NewRuleEnv("ENV_2", "^abc[0-9]+")
	if r2.Match(dv1) {
		t.Fatal("The device `dv1` was found incorrectly")
	}

	r3 := NewRuleEnv("ENV_2", "test data")
	if !r3.Match(dv1) {
		t.Fatal("Could not find device `dv1`")
	}

	r4 := NewRuleEnv("ENV_NOT_EXIST", "test data")
	if r4.Match(dv1) {
		t.Fatal("The device `dv1` was found incorrectly")
	}

	r5 := RuleEnv{regexp: nil, envName: "ATTR_1"}
	if r5.Match(dv1) {
		t.Fatal("The device `dv1` was found incorrectly")
	}
}
