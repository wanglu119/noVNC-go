// +build !windows

package utils

import (
	"os"
	"syscall"
)

/**
* @Description:
* @param path
* @return int64
 */
func GetFileCreateTime(path string) int64 {
	fileInfo, _ := os.Stat(path)

	stat_t := fileInfo.Sys().(*syscall.Stat_t)
	tCreate := int64(stat_t.Ctim.Sec)
	return tCreate
}
