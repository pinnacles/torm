package torm

import (
	"testing"

	"github.com/pinnacles/torm/internal/test"
)

func TestRegister(t *testing.T) {
	Register(test.TestSchema{})

	if m, ok := metas["test"]; !ok {
		t.Fatal("metas don't have a `test` key")
	} else {
		if m.TableName != "test" {
			t.Error("m.TableName is not `test`")
		}
		ok := false
		for _, f := range m.Fields {
			if f == "foo" {
				ok = true
			}
		}
		if !ok {
			t.Error("m.Fileds don't have `foo` key")
		}
		ok = false
		for _, f := range m.Fields {
			if f == "bar" {
				ok = true
			}
		}
		if ok {
			t.Error("m.Fileds have `bar` key")
		}
	}
}

func TestVerboseLevelValid(t *testing.T) {
	for _, lv := range []int{0, 1, 2, 3} {
		VerboseLevel(lv)
	}
}

func TestVerboseLevelInvalidNegative(t *testing.T) {
	lv := -1

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("VerboseLevel(%d) don't panic", lv)
		} else if err != "VerboseLevel is must be in range of 0 to 3" {
			t.Errorf("VerboseLevel(%d) error message is not that was expected", lv)
		}
	}()

	VerboseLevel(lv)
}

func TestVerboseLevelInvalidPositive(t *testing.T) {
	lv := 4

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("VerboseLevel(%d) don't panic", lv)
		} else if err != "VerboseLevel is must be in range of 0 to 3" {
			t.Errorf("VerboseLevel(%d) error message is not that was expected", lv)
		}
	}()

	VerboseLevel(lv)
}
