package utlis

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func GetCurrentProcessID() string{
	return strconv.Itoa(os.Getegid())
}

func GetCurrentGoroutineID()string{
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf,false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo,"[running]")[0], "goroutine")[1])
}

func GetProcessAndGoroutineIDStr()string{
	return fmt.Sprintf("%s_%s", GetCurrentProcessID(), GetCurrentGoroutineID())
}
