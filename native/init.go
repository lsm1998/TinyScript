package native

import (
	"errors"
	"reflect"
)

var NativeObjectNotFindErr = errors.New("native object not find")
var NativeMethodNotFindErr = errors.New("native method not find")

type NativeFunc struct {
	Func   reflect.Value
	Params []string
	IsErr  bool
}

var nativeMap = map[string]map[string]*NativeFunc{}

func registerNative(name string, method string, fun *NativeFunc) {
	objMap := nativeMap[name]
	if objMap == nil {
		objMap = map[string]*NativeFunc{}
	}
	objMap[method] = fun
	nativeMap[name] = objMap
}

func GetNativeObject(name string) map[string]*NativeFunc {
	return nativeMap[name]
}

func CallNativeMethod(obj string, method string, args ...interface{}) (interface{}, error) {
	instance, ok := nativeMap[obj]
	if !ok {
		return nil, NativeObjectNotFindErr
	}

	m := instance[method]
	if !m.Func.IsValid() {
		return nil, NativeMethodNotFindErr
	}

	var values = make([]reflect.Value, 0, len(args))
	for _, v := range args {
		values = append(values, reflect.ValueOf(v))
	}
	callResult := m.Func.Call(values)

	var result = make([]interface{}, 0, len(callResult))

	for _, v := range callResult {
		result = append(result, v.Interface())
	}
	return result, nil
}
