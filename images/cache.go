package images

import (
	"bytes"

	"github.com/adm87/stellar/engine/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var Cache = assets.NewAssetCache[ebiten.Image](1024,
	// AssetAllocator function that creates an ebiten.Image from raw byte data
	func(data []byte) (*ebiten.Image, error) {
		img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))

		if err != nil {
			return nil, err
		}

		return img, nil
	},

	// AssetDeallocator function that deallocates an ebiten.Image when it is removed from the cache
	func(img *ebiten.Image) {
		img.Deallocate()
	},
)
