package client

import "runtime"

func CallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
			return ""
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
			return ""
	}
	return f.Name()
}