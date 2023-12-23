package native

import (
	"os"
	"reflect"
	"runtime"
)

type nativeOs struct {
}

func (*nativeOs) NumCPU() int {
	return runtime.NumCPU()
}

func (*nativeOs) Pwd() string {
	dir, _ := os.Getwd()
	return dir
}

func init() {
	registerNative("os", "NumCPU", &NativeFunc{
		Func:   reflect.ValueOf(new(nativeOs)).MethodByName("NumCPU"),
		Params: nil,
	})
	registerNative("os", "Pwd", &NativeFunc{
		Func:   reflect.ValueOf(new(nativeOs)).MethodByName("Pwd"),
		Params: nil,
	})

}
