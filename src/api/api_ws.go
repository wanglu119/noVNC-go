package api

import (
	"net"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/wanglu119/me-deps/webCommon"
	"github.com/wanglu119/me-deps/webCommon/ws"
)

// use default options
// var upgrader = websocket.Upgrader{}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *noVNCApi) vncConnect(w http.ResponseWriter, r *http.Request, d webCommon.WebData) {
	remoteConn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error("upgrade: ", err)
		return
	}
	defer remoteConn.Close()
	wswConn := &ws.WsWrapper{Conn: remoteConn}

	ClientAddr := r.URL.Query().Get("vncServer")
	// start up the private connection
	localConn, err := net.Dial("tcp", ClientAddr)
	if err != nil {
		log.Error("Failed to open private log ", ClientAddr, err)
		return
	} else {
		log.Info(r.RemoteAddr, " connect vnc server: ", ClientAddr)
	}
	defer localConn.Close()

	errs := make(chan error, 2)

	go func() {
		errs <- func() error {
			for {
				data, err := wswConn.Read()
				if err != nil {
					log.Error(err)
					return err
				}

				_, err = localConn.Write(data)
				if err != nil {
					log.Error(err)
					return err
				}
			}
		}()
	}()

	go func() {
		errs <- func() error {
			for {
				data := make([]byte, 1024)

				n, err := localConn.Read(data)
				if err != nil {
					log.Error(err)
					return err
				}

				err = wswConn.Write(data[:n])
				if err != nil {
					log.Error(err)
					return err
				}
			}
		}()
	}()

	select {
	case err = <-errs:
	}
	log.Info(r.RemoteAddr, " disconnect vnc server: ", ClientAddr)
}
