package webCommon

import (
	"net/http"

	"github.com/spf13/afero"
)

type WebData interface {
	GetResponse() *Response
	SetResponse(*Response)
	GetAuthToken() string
	GetAuthData() interface{}
	SetAuthData(interface{})
	SetHttpHeader(http.Header)

	GetFs() afero.Fs
}

type WebDataImplPart struct {
	Resp   *Response
	Header http.Header
}

func (wdip *WebDataImplPart) GetResponse() *Response {
	return wdip.Resp
}

func (wdip *WebDataImplPart) SetResponse(resp *Response) {
	wdip.Resp = resp
}
func (wdip *WebDataImplPart) SetHttpHeader(header http.Header) {
	wdip.Header = header
}

type HandleFunc func(w http.ResponseWriter, r *http.Request, d WebData)
type CreateWebDataFunc func() WebData
type WebHandler struct {
	CreateWebData CreateWebDataFunc
}

func (wh *WebHandler) handle(fn HandleFunc, prefix string) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := wh.CreateWebData()
		res := NewResponse(w)
		d.SetResponse(res)
		d.SetHttpHeader(r.Header)
		defer res.SendJson()
		fn(w, r, d)
	})

	return http.StripPrefix(prefix, handler)
}

func (wh *WebHandler) Monkey(fn HandleFunc, prefix string) http.Handler {
	return wh.handle(fn, prefix)
}
