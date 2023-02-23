package api

import (
	"bytes"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/dchest/captcha"

	"github.com/wanglu119/me-deps/webCommon"
)

func (ca *commonApi) captchaId(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
	res := d.GetResponse()
	res.Data = captcha.New()
}

func (ca *commonApi) captchaImage(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
	res := d.GetResponse()
	defer func() {
		if res.Status == http.StatusOK {
			res.IsSend = false
		}
	}()

	imgWidth := captcha.StdWidth
	imgHeight := captcha.StdHeight

	serve := func(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool) error {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		var content bytes.Buffer
		var err error
		switch ext {
		case ".png":
			w.Header().Set("Content-Type", "image/png")
			err = captcha.WriteImage(&content, id, imgWidth, imgHeight)
		case ".wav":
			w.Header().Set("Content-Type", "audio/x-wav")
			err = captcha.WriteAudio(&content, id, lang)
		default:
			return captcha.ErrNotFound
		}

		if err != nil {
			return err
		}

		if download {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
		http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
		return nil
	}

	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		errMsg := "ext or id empty"
		log.Error(errMsg)
		res.Data = errMsg
		res.Status = http.StatusNotFound
		return
	}
	if r.FormValue("reload") != "" {
		captcha.Reload(id)
	}
	lang := strings.ToLower(r.FormValue("lang"))
	download := path.Base(dir) == "download"
	if serve(w, r, id, ext, lang, download) == captcha.ErrNotFound {
		errMsg := fmt.Sprintf("not found: %s%s", id, ext)
		log.Error(errMsg)
		res.Data = errMsg
		res.Status = http.StatusNotFound
		return
	}
}
