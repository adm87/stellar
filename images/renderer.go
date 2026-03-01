package images

import (
	"github.com/adm87/stellar/engine/structures/store"
	"github.com/hajimehoshi/ebiten/v2"
)

func RenderImage(screen *ebiten.Image, id store.StoreID, opt *ebiten.DrawImageOptions) {
	img, ok := imageStore.Get(id)

	if !ok {
		return
	}

	screen.DrawImage(img, opt)
}
