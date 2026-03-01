package assets

import (
	"strings"
	"sync"

	"github.com/adm87/stellar/engine/structures/store"
)

// AssetPath represents a game asset path. The path is relative to the filesystem it is loaded from.
type AssetPath string

// String returns the string representation of the Asset, which is its path. This is useful for logging and debugging purposes.
func (a AssetPath) String() string {
	return string(a)
}

// Type returns the AssetType of the Asset by extracting the file extension from its path.
// For example, if the asset is "images/sprite.png", it will return "png". This is used to determine how to import the asset based on its type.
func (a AssetPath) Type() AssetType {
	ext := strings.LastIndex(string(a), ".")
	if ext == -1 || ext < strings.LastIndex(string(a), "/") || ext == len(string(a))-1 {
		return ""
	}
	return AssetType(string(a)[ext+1:])
}

// AssetType represents a type of asset, such as "jpeg", "png", "json", etc. This can be used to determine how to load the asset.
type AssetType string

// AssetImporter is an interface that defines how to import an asset.
// It has a single method, Import, which takes an AssetPath and its raw byte data, and returns an error if the import fails.
type AssetImporter interface {
	Import(path AssetPath, rawData []byte) error
}

// AssetAllocator is a function type that defines how to allocate an asset.
// It takes the raw byte data of the asset and returns a pointer to the allocated asset and an error if the allocation fails.
type AssetAllocator[T any] func(rawData []byte) (*T, error)

// AssetDeallocator is a function type that defines how to deallocate an asset.
// It takes a pointer to the asset that needs to be deallocated.
type AssetDeallocator[T any] func(*T)

// AssetCache is a generic structure that manages the caching of loaded assets.
// It uses a store to manage the actual asset data and a registry to map asset paths to their corresponding store IDs.
type AssetCache[T any] struct {
	cache     *store.Store[*T]
	allocFn   AssetAllocator[T]
	deallocFn AssetDeallocator[T]

	assets map[AssetPath]store.StoreID
	mu     sync.RWMutex
}

// NewAssetCache creates a new AssetCache with the specified capacity for the underlying store.
// The capacity determines how many assets can be stored before the cache needs to evict old assets.
func NewAssetCache[T any](capacity uint32, allocFn AssetAllocator[T], deallocFn AssetDeallocator[T]) *AssetCache[T] {
	return &AssetCache[T]{
		cache:     store.NewStore[*T](capacity),
		allocFn:   allocFn,
		deallocFn: deallocFn,
		assets:    make(map[AssetPath]store.StoreID),
	}
}

// Import imports an asset into the cache using the provided AssetPath and raw byte data.
// It uses the AssetAllocator function to allocate the asset and stores it in the cache.
func (ac *AssetCache[T]) Import(path AssetPath, rawData []byte) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.assets[path]; exists {
		return DuplicateAsset{Asset: path}
	}

	item, err := ac.allocFn(rawData)

	if err != nil {
		return err
	}

	id, err := ac.cache.Allocate(item)

	if err != nil {
		ac.deallocFn(item)
		return err
	}

	ac.assets[path] = id
	return nil
}

// GetID retrieves the StoreID of an asset from the cache based on its AssetPath. It returns the StoreID and a boolean
// indicating whether the asset was found in the cache.
func (ac *AssetCache[T]) GetID(asset AssetPath) (store.StoreID, bool) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	id, exists := ac.assets[asset]
	return id, exists
}

// GetByID retrieves an asset from the cache based on its StoreID. It returns a pointer to the asset data and a boolean
// indicating whether the asset was found in the cache.
func (ac *AssetCache[T]) GetByID(id store.StoreID) (*T, bool) {
	return ac.cache.Get(id)
}

// Add adds a new asset to the cache. It takes an AssetPath and a pointer to the asset data. If the asset already exists in the
// cache, it returns a DuplicateAsset error.
func (ac *AssetCache[T]) Add(path AssetPath, item *T) (store.StoreID, error) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.assets[path]; exists {
		var zero store.StoreID
		return zero, DuplicateAsset{Asset: path}
	}

	id, err := ac.cache.Allocate(item)

	if err != nil {
		var zero store.StoreID
		return zero, err
	}

	ac.assets[path] = id
	return id, nil
}

// Remove removes an asset from the cache based on its AssetPath.
// Calls the provided deallocation function to free the memory used by the asset before removing it from the cache.
func (ac *AssetCache[T]) Remove(path AssetPath) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if _, exists := ac.assets[path]; !exists {
		return
	}

	item, exists := ac.cache.Get(ac.assets[path])

	if exists {
		ac.deallocFn(item)
	}

	ac.cache.Deallocate(ac.assets[path])
	delete(ac.assets, path)
}

// Clear clears all assets from the cache.
// Calls the provided deallocation function for each asset to free their memory before clearing the cache.
func (ac *AssetCache[T]) Clear() {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	for _, id := range ac.assets {
		item, exists := ac.cache.Get(id)

		if exists {
			ac.deallocFn(item)
		}

		ac.cache.Deallocate(id)
	}

	ac.assets = make(map[AssetPath]store.StoreID)
}
