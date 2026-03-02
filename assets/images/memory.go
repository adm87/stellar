package images

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func AllocateImage(data []byte) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	return img, nil
}

func DeallocateImage(img *ebiten.Image) {
	img.Deallocate()
}
