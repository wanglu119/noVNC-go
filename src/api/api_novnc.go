package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/rs/xid"

	"github.com/wanglu119/me-deps/webCommon"
)

type noVnc struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

var noVncMap sync.Map

func (s *noVNCApi) createNoVnc(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
	res := d.GetResponse()
	var err error
	defer func() {
		if err != nil {
			webCommon.ProcError(res, err)
			return
		}
	}()

	param := &noVnc{}
	err = json.NewDecoder(r.Body).Decode(param)
	if err != nil {
		log.Error(err)
		return
	}

	guid := xid.New()
	param.Id = guid.String()

	noVncMap.Store(param.Id, param)
}

func (s *noVNCApi) listNoVnc(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
	res := d.GetResponse()
	var err error
	defer func() {
		if err != nil {
			webCommon.ProcError(res, err)
			return
		}
	}()

	arr := []*noVnc{}

	noVncMap.Range(func(k any, v any) bool {
		arr = append(arr, v.(*noVnc))
		return true
	})

	res.Data = arr
}

func (s *noVNCApi) deleteNoVnc(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
	res := d.GetResponse()
	var err error
	defer func() {
		if err != nil {
			webCommon.ProcError(res, err)
			return
		}
	}()

	param := &noVnc{}
	err = json.NewDecoder(r.Body).Decode(param)
	if err != nil {
		log.Error(err)
		return
	}

	noVncMap.Delete(param.Id)
}
