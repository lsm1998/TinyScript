package native

import (
	"testing"
)

func TestCallNativeMethod(t *testing.T) {
	t.Log(CallNativeMethod("file", "ReadFile", "init_test.go"))
}
