package images

import (
	"github.com/adm87/stellar/engine/assets"
	"github.com/adm87/stellar/engine/structures/store"
	"github.com/hajimehoshi/ebiten/v2"
)

func deallocateImage(img *ebiten.Image) {
	img.Deallocate()
}

var imageCache = assets.NewAssetCache[ebiten.Image](1024)

func GetImageID(path assets.AssetPath) (store.StoreID, bool) {
	return imageCache.GetID(path)
}

func GetImage(id store.StoreID) (*ebiten.Image, bool) {
	return imageCache.GetByID(id)
}

func AddImage(path assets.AssetPath, img *ebiten.Image) error {
	_, err := imageCache.Add(path, img)
	return err
}

func RemoveImage(path assets.AssetPath) {
	imageCache.Remove(path, deallocateImage)
}

func ClearCache() {
	imageCache.Clear(deallocateImage)
}
