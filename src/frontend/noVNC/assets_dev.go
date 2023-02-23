//go:build debug
// +build debug

package noVNC

import (
	"io/fs"
	"os"
)

var assets fs.FS = os.DirFS("/data/Project/GitLab/WL_Ngrok_V2/src/frontend/noVNC")

func Assets() fs.FS {
	return assets
}
