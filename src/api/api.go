package api

import (
	"io/fs"

	noVNCFrontend "github.com/wanglu119/noVNC-go/frontend/noVNC"

	Log "github.com/wanglu119/me-deps/log"
	"github.com/wanglu119/me-deps/mux"
	"github.com/wanglu119/me-deps/webCommon"
	CommonApi "github.com/wanglu119/me-deps/webCommon/api"
)

var log = Log.GetLogger()

type noVNCApi struct {
	*webCommon.WebHandler
	*webCommon.JwtAuth
}

func NewNoVNCApi() *noVNCApi {
	return &noVNCApi{
		WebHandler: &webCommon.WebHandler{CreateWebData: CreateWebData},
		JwtAuth:    &webCommon.JwtAuth{},
	}
}

func (s *noVNCApi) Serve(router *mux.Router, port uint32) {
	api := router.PathPrefix("/api/ng-noVNC").Subrouter()

	api.Handle("/websockify", s.Monkey(s.vncConnect, "")).Methods("OPTIONS", "GET")
	api.Handle("/create", s.Monkey(s.createNoVnc, "")).Methods("OPTIONS", "POST")
	api.Handle("/list", s.Monkey(s.listNoVnc, "")).Methods("OPTIONS", "GET")
	api.Handle("/delete", s.Monkey(s.deleteNoVnc, "")).Methods("OPTIONS", "POST")

	// static api
	assetsFs, err := fs.Sub(noVNCFrontend.Assets(), ".")
	if err != nil {
		panic(err)
	}
	index, static := CommonApi.GetStaticHandlers(assetsFs)
	router.PathPrefix("/").Handler(s.Monkey(static, "/"))
	router.NotFoundHandler = s.Monkey(index, "")

	mux.Serve(router, port)
}

func (s *noVNCApi) Setup(router *mux.Router) {
	api := router.PathPrefix("/api/ng-noVNC").Subrouter()

	api.Handle("/websockify", s.Monkey(s.WithJwt(s.vncConnect), "")).Methods("OPTIONS", "GET")

	// static api
	assetsFs, err := fs.Sub(noVNCFrontend.Assets(), ".")
	if err != nil {
		panic(err)
	}
	index, static := CommonApi.GetStaticHandlers(assetsFs)
	router.PathPrefix("/noVNC/").Handler(s.Monkey(static, "/noVNC/"))
	router.PathPrefix("/").Handler(s.Monkey(static, "/"))
	router.NotFoundHandler = s.Monkey(index, "")
}
