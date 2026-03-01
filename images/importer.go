package images

import (
	"bytes"

	"github.com/adm87/stellar/engine/assets"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func importImage(asset assets.AssetPath, rawData []byte) error {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewBuffer(rawData))

	if err != nil {
		return err
	}

	if err := AddImage(asset, img); err != nil {
		deallocateImage(img)
		return err
	}

	return nil
}
