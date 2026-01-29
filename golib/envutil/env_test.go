package envutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsProd(t *testing.T) {
	_ = os.Setenv("APP_ENV", Prod)
	if !IsProd() {
		t.Error("IsProd() should return true")
	}
}

func TestIsTest(t *testing.T) {
	_ = os.Setenv("APP_ENV", Test)
	if !IsTest() {
		t.Error("IsTest() should return true")
	}
}

func TestIsDev(t *testing.T) {
	_ = os.Setenv("APP_ENV", Dev)
	if !IsDev() {
		t.Error("IsDev() should return true")
	}
}

func TestEnv(t *testing.T) {
	_ = os.Setenv("APP_ENV", Dev)
	assert.Equal(t, "dev", Env())
}
