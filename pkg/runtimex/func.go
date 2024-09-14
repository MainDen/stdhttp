package runtimex

import (
	"reflect"
	"runtime"
)

func GetFuncName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}
