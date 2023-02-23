package utils

import (
	"os"
	"path"
)

func GetPWD() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func GetPwdWith(sub string) string {
	return path.Join(GetPWD(), sub)
}
