package images

import (
	"bytes"

	"github.com/adm87/stellar/engine/assets"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func importImage(asset assets.Asset, rawData []byte) error {
	imageStoreMu.Lock()
	defer imageStoreMu.Unlock()

	if _, exists := imageRegistry[asset]; exists {
		return assets.DuplicateAsset{Asset: asset}
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewBuffer(rawData))

	if err != nil {
		return err
	}

	id, err := imageCache.Allocate(img)

	if err != nil {
		return err
	}

	imageRegistry[asset] = id
	return nil
}
