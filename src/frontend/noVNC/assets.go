//go:build !debug
// +build !debug

package noVNC

import "embed"

//go:embed *
var assets embed.FS

func Assets() embed.FS {
	return assets
}
