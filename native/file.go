package native

import (
	"io/fs"
	"os"
	"reflect"
)

type nativeFile struct {
}

func (*nativeFile) ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	return string(content), err
}

func (n *nativeFile) WriteFile(filename string, content string, mode string) error {
	flag := os.O_CREATE
	if mode == "append" {
		flag |= os.O_APPEND
	}
	file, err := os.OpenFile(filename, flag, fs.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
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
