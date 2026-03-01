package images

import (
	"sync"

	"github.com/adm87/stellar/engine/assets"
)

var registerOnce sync.Once

// Register registers the image asset importers for supported image types.
func Register() {
	registerOnce.Do(func() {
		assets.RegisterImporter("jpg", Cache)
		assets.RegisterImporter("jpeg", Cache)
		assets.RegisterImporter("png", Cache)
	})
}
