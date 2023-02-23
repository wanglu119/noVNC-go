package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/wenlng/go-captcha/captcha"

	"github.com/wanglu119/me-deps/utils"
	"github.com/wanglu119/me-deps/webCommon"
)

func (ca *commonApi) getCaptchaData(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
	res := webCommon.NewResponse(w)
	defer res.SendJson()

	capt := captcha.GetCaptcha()

	dots, b64, tb64, key, err := capt.Generate()
	if err != nil {
		webCommon.ProcError(res, err)
		log.Error(err)
		return
	}
	fmt.Println(dots)
	writeCache(dots, key)
	fmt.Println(utils.GetPWD())
	res.Data = map[string]interface{}{
		"code":         0,
		"image_base64": b64,
		"thumb_base64": tb64,
		"captcha_key":  key,
	}
}

func writeCache(v interface{}, file string) {
	bt, _ := json.Marshal(v)
	month := time.Now().Month().String()
	cacheDir := path.Join(getCacheDir(), month)
	_ = os.MkdirAll(cacheDir, 0660)
	file = path.Join(cacheDir, file+".json")
	logFile, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer logFile.Close()
	checkCacheOvertimeFile()
	_, _ = io.WriteString(logFile, string(bt))
}

func readCache(file string) string {
	month := time.Now().Month().String()
	cacheDir := getCacheDir() + month + "/"
	file = cacheDir + file + ".json"

	if !checkFileIsExist(file) {
		return ""
	}

	bt, err := ioutil.ReadFile(file)
	err = os.Remove(file)
	if err == nil {
		return string(bt)
	}
	return ""
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func getCacheDir() string {
	return utils.GetPwdWith("bin/.cache/captcha")
}

func checkCacheOvertimeFile() {
	files, files1, _ := listDir(getCacheDir())
	for _, table := range files1 {
		temp, _, _ := listDir(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	for _, file := range files {
		t := utils.GetFileCreateTime(file)
		ex := time.Now().Unix() - t
		if ex > (60 * 30) {
			_ = os.Remove(file)
		}
	}
}

func listDir(dirPth string) (files []string, files1 []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { //
			files1 = append(files1, dirPth+PthSep+fi.Name())
			_, _, _ = listDir(dirPth + PthSep + fi.Name())
		} else {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, files1, nil
}
