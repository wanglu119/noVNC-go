package mux

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type Router = mux.Router

func NewRouter() *Router {
	return mux.NewRouter()
}

func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func Serve(router *Router, port uint32) {
	var listener net.Listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Listening on: ", port)
	err = http.Serve(listener, router)
	if err != nil {
		panic(err)
	}
}
