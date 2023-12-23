package native

import (
	"reflect"
	"time"
)

type nativeTime struct {
}

func (*nativeTime) Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (*nativeTime) Timestamp() int64 {
	return time.Now().UnixMilli()
}

func init() {
	registerNative("time", "Time", &NativeFunc{
		Func: reflect.ValueOf(new(nativeTime)).MethodByName("Time"),
	})
	registerNative("time", "Timestamp", &NativeFunc{
		Func: reflect.ValueOf(new(nativeTime)).MethodByName("Timestamp"),
	})
}
