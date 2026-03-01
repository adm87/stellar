package images

import (
	"sync"

	"github.com/adm87/stellar/engine/assets"
	"github.com/adm87/stellar/engine/structures/store"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	imageCache    = store.NewStore[*ebiten.Image](1024)
	imageRegistry = make(map[assets.Asset]store.StoreID)
	imageStoreMu  sync.RWMutex
)

// GetStoreID retrieves the StoreID associated with the given Asset.
// It returns the StoreID and a boolean indicating whether the Asset was found in the registry.
func GetStoreID(asset assets.Asset) (store.StoreID, bool) {
	imageStoreMu.RLock()
	defer imageStoreMu.RUnlock()

	id, exists := imageRegistry[asset]

	return id, exists
}

// GetByAsset retrieves an image from the store using its associated Asset. It returns the image and a boolean indicating whether the image was found.
func GetByAsset(asset assets.Asset) (*ebiten.Image, bool) {
	imageStoreMu.RLock()
	defer imageStoreMu.RUnlock()

	id, exists := imageRegistry[asset]

	if !exists {
		return nil, false
	}

	img, ok := imageCache.Get(id)

	if !ok {
		return nil, false
	}

	return img, true
}

// GetByID retrieves an image from the store using its StoreID. It returns the image and a boolean indicating whether the image was found.
func GetByID(id store.StoreID) (*ebiten.Image, bool) {
	img, ok := imageCache.Get(id)

	if !ok {
		return nil, false
	}

	return img, true
}

// Add adds a new image to the store and returns its StoreID. If the asset already exists in the registry, it returns an error.
func Add(asset assets.Asset, img *ebiten.Image) (store.StoreID, error) {
	imageStoreMu.Lock()
	defer imageStoreMu.Unlock()

	if _, exists := imageRegistry[asset]; exists {
		return store.StoreID{}, assets.DuplicateAsset{Asset: asset}
	}

	id, err := imageCache.Allocate(img)

	if err != nil {
		return store.StoreID{}, err
	}

	imageRegistry[asset] = id
	return id, nil
}

// Remove removes an image from the store and registry based on the given Asset. If the asset does not exist in the registry, it does nothing.
//
// Note: The resource is released from memory and should not be used.
func Remove(asset assets.Asset) {
	imageStoreMu.Lock()
	defer imageStoreMu.Unlock()

	id, exists := imageRegistry[asset]

	if !exists {
		return
	}

	img, ok := imageCache.Get(id)

	if ok {
		img.Deallocate()
	}

	imageCache.Deallocate(id)

	delete(imageRegistry, asset)
}

// Clear removes all images from the store and registry.
//
// Note: The resources are released from memory and should not be used.
func Clear() {
	imageStoreMu.Lock()
	defer imageStoreMu.Unlock()

	for _, id := range imageRegistry {
		img, ok := imageCache.Get(id)

		if ok {
			img.Deallocate()
		}

		imageCache.Deallocate(id)
	}

	imageRegistry = make(map[assets.Asset]store.StoreID)
}
