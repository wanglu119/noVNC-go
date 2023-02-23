package main

import (
	"github.com/spf13/cobra"

	Log "github.com/wanglu119/me-deps/log"
	"github.com/wanglu119/me-deps/mux"
	"github.com/wanglu119/me-deps/webCommon"
	"github.com/wanglu119/noVNC-go/api"
)

var log = Log.GetLogger()
var apiPort uint32 = 0

var rootCmd = &cobra.Command{
	Use:   "noVNC-go",
	Short: "noVNC-go",
	Run: func(cmd *cobra.Command, args []string) {

		router := mux.NewRouter()
		router.Use(webCommon.CorsMiddleware)
		serverApi := api.NewNoVNCApi()
		serverApi.Serve(router, apiPort)
	},
}

func init() {
	flags := rootCmd.Flags()

	flags.Uint32Var(&apiPort, "api-port", 9002, "api port")
}

func main() {
	rootCmd.Execute()
}
