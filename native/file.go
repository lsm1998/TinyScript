package native

import (
	"os"
	"reflect"
)

type nativeFile struct {
}

func (*nativeFile) ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	return string(content), err
}

func (*nativeFile) WriteFile(filename string, content string, mode int32) error {
	return os.WriteFile(filename, []byte(content), os.FileMode(mode))
}

func init() {
	registerNative("file", "ReadFile", &NativeFunc{
		Func:   reflect.ValueOf(new(nativeFile)).MethodByName("ReadFile"),
		Params: []string{"filename"},
		IsErr:  true,
	})
	registerNative("file", "WriteFile", &NativeFunc{
		Func:   reflect.ValueOf(new(nativeFile)).MethodByName("WriteFile"),
		Params: []string{"filename", "content", "mode"},
		IsErr:  true,
	})
}
