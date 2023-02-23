package api

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/wanglu119/me-deps/webCommon"
)

func handleWithStaticData(w http.ResponseWriter, _ *http.Request, d webCommon.WebData, fSys fs.FS, file, contentType string) {
	w.Header().Set("Content-Type", contentType)
	res := d.GetResponse()

	fileContents, err := fs.ReadFile(fSys, file)
	if err != nil {
		if err == os.ErrNotExist {
			res.Status = http.StatusNotFound
			res.SendJson()
			return
		}
		res.Status = http.StatusNotFound
		res.SendJson()
		res.Data = fmt.Sprintf("%v", err)
		return
	}
	res.IsSend = false
	w.Write(fileContents)

	return
}

func GetStaticHandlers(assetsFs fs.FS) (index, static webCommon.HandleFunc) {
	index = func(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
		res := d.GetResponse()
		if r.Method != http.MethodGet {
			res.Status = http.StatusNotFound
			res.SendJson()
			return
		}
		res.IsSend = false

		w.Header().Set("x-xss-protection", "1; mode=block")
		handleWithStaticData(w, r, d, assetsFs, "index.html", "text/html; charset=utf-8")
	}

	static = func(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
		res := d.GetResponse()
		if r.Method != http.MethodGet {
			res.Status = http.StatusNotFound
			res.SendJson()
			return
		}
		var err error
		defer func() {
			if err == nil {
				res.IsSend = false
			}
		}()

		const maxAge = 86400 // 1 day
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v", maxAge))

		urlPath := r.URL.Path
		if strings.Count(urlPath, ".") == 0 {
			urlPath = "index.html"
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		} else {
			if strings.HasSuffix(urlPath, ".js") {
				w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			} else if strings.HasSuffix(urlPath, ".css") {
				w.Header().Set("Content-Type", "text/css; charset=utf-8")
			} else {
				http.FileServer(http.FS(assetsFs)).ServeHTTP(w, r)
				return
			}
			if _, err := fs.Stat(assetsFs, urlPath+".gz"); err == nil {
				urlPath = urlPath + ".gz"
				w.Header().Set("Content-Encoding", "gzip")
			}
		}

		fileContents, err := fs.ReadFile(assetsFs, urlPath)
		if err != nil {
			res.Status = http.StatusNotFound
			res.Data = fmt.Sprintf("%v", err)
			log.Error(err)
			res.SendJson()
			return
		}

		if _, err = w.Write(fileContents); err != nil {
			res.Status = http.StatusInternalServerError
			res.Data = fmt.Sprintf("%v", err)
			res.SendJson()
			return
		}
		res.IsSend = false

		return
	}

	return index, static
}
