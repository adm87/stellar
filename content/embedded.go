package content

import (
	"embed"

	"github.com/adm87/stellar/assets"
)

//go:embed embedded
var EmbeddedFS embed.FS

const (
	EmbeddedImage10x10 assets.AssetPath = "embedded/images/image_10x10.png"
)
