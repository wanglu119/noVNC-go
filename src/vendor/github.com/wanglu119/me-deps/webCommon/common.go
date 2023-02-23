package webCommon

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	w      http.ResponseWriter
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	IsSend bool        `json:"-"`
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		w:      w,
		Status: 200,
		Data:   nil,
		IsSend: true,
	}
}

func (r *Response) SendJson() (int, error) {
	if r.IsSend {
		resp, err := json.Marshal(r)
		if err != nil {
			log.Error(err)
		}
		r.w.Header().Set("Content-Type", "application/json")
		r.w.WriteHeader(r.Status)
		return r.w.Write(resp)
	} else {
		return 0, nil
	}
}

func ProcError(res *Response, err error) {
	errMsg := fmt.Sprintf("%v", err)
	res.Data = errMsg
	res.Status = http.StatusInternalServerError
}

// -----------------------------------------------------------

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "close")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT,PATCH")
		//w.Header().Set("Access-Control-Allow-Headers", `Origin,X-Requested-With,Content-Type,Accept,Authorization,token,X-Auth,x-auth,captchaId,captchaVal`)
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// w.Header().Set("Strict-Transport-Security","max-age=0; includeSubDomains")

		if r.Method == http.MethodOptions {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			return
		}
		next.ServeHTTP(w, r)
	})
}
