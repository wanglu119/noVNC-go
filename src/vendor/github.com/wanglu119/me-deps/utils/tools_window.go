// +build windows

package utils

import (
	"os"
	"runtime"
	"syscall"
	"time"
)

/**
* @Description:
* @param path
* @return int64
 */
func GetFileCreateTime(path string) int64 {
	osType := runtime.GOOS
	fileInfo, _ := os.Stat(path)
	if osType == "windows" {
		wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
		tNanSeconds := wFileSys.CreationTime.Nanoseconds() /// ns
		tSec := tNanSeconds / 1e9                          /// s
		return tSec
	}

	return time.Now().Unix()
}
