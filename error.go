package gutils

import (
	"fmt"
	"runtime"

	"github.com/yinyihanbing/gutils/logs"
)

// 崩溃错误处理
func PanicError() error {
	if r := recover(); r != nil {
		buf := make([]byte, 4096)
		l := runtime.Stack(buf, false)
		err := fmt.Errorf("%v: %s", r, buf[:l])
		logs.Error(err)
		return err
	}
	return nil
}
