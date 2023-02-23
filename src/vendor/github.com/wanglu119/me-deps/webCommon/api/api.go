package api

import (
	"github.com/wanglu119/me-deps/mux"
	"github.com/wanglu119/me-deps/webCommon"
)

type commonApi struct {
	*webCommon.WebHandler
}

func NewCommonApi(createWebData webCommon.CreateWebDataFunc) *commonApi {
	return &commonApi{
		&webCommon.WebHandler{
			CreateWebData: createWebData,
		},
	}
}

func (ca *commonApi) Setup(router *mux.Router) {
	common := router.PathPrefix("/api/common").Subrouter()
	common.Handle("/captcha/id", ca.Monkey(ca.captchaId, "")).Methods("OPTIONS", "GET")
	common.Handle("/captcha/image/{id}", ca.Monkey(ca.captchaImage, "/api/common/captcha/image")).Methods("OPTIONS", "GET")

	common.Handle("/captcha_v2/getData", ca.Monkey(ca.getCaptchaData, "")).Methods("OPTIONS", "GET")
}
