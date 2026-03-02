package assets

import (
	"io"
	"io/fs"
	"strings"
	"sync"

	"github.com/adm87/stellar/assets/images"
	"github.com/adm87/stellar/errs"
	"github.com/hajimehoshi/ebiten/v2"
)

// AssetType is the file extension describing the type of asset (e.g., "png", "json", "txt").
type AssetType string

const InvalidAssetType AssetType = ""

// AssetPath is the relative path to an asset file within an asset filesystem.
type AssetPath string

// Type returns the AssetType based on the file extension of the AssetPath.
func (p AssetPath) Type() AssetType {
	if idx := strings.LastIndexByte(string(p), '.'); idx != -1 {
		return AssetType(p[idx+1:])
	}
	return InvalidAssetType
}

// AllocateAssetFunc is a function type for allocating an asset of type T from raw byte data.
type AllocateAssetFunc[T any] func(data []byte) (T, error)

// DeallocateAssetFunc is a function type for deallocating an asset of type T, allowing for any necessary cleanup.
type DeallocateAssetFunc[T any] func(asset T)

// AssetManifest maps AssetPaths to their corresponding AssetIDs in the AssetStore.
type AssetManifest map[AssetPath]AssetID

// AssetProcessor defines the interface for allocating and deallocating assets of a specific type, given their AssetPath and raw byte data.
type AssetProcessor interface {
	allocate(path AssetPath, reader io.Reader) error
	deallocate(path AssetPath) error
}

// AssetCache manages a cache of assets of type T, using an AssetStore for storage and providing allocation and deallocation functions.
type AssetCache[T any] struct {
	store    *AssetStore[T]
	alloc    AllocateAssetFunc[T]
	dealloc  DeallocateAssetFunc[T]
	manifest AssetManifest
	mu       sync.RWMutex
}

// allocate loads an asset from raw byte data, allocating it using the provided AllocateAssetFunc and storing it in the AssetStore.
// If the asset is already loaded, it does nothing. Returns an error if allocation or storage fails.
func (c *AssetCache[T]) allocate(path AssetPath, reader io.Reader) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.manifest[path]; exists {
		return nil
	}

	data, err := io.ReadAll(reader)

	if err != nil {
		return errs.InvalidOperation{
			Message: "failed to read asset data: " + err.Error(),
		}
	}

	asset, err := c.alloc(data)

	if err != nil {
		return errs.InvalidOperation{
			Message: "failed to allocate asset: " + err.Error(),
		}
	}

	id, err := c.store.Add(asset)

	if err != nil {
		c.dealloc(asset)
		return errs.InvalidOperation{
			Message: "failed to add asset to store: " + err.Error(),
		}
	}

	c.manifest[path] = id
	return nil
}

// deallocate removes an asset from the cache and deallocates it using the provided DeallocateAssetFunc.
// If the asset is not found or removal fails, returns an error.
func (c *AssetCache[T]) deallocate(path AssetPath) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	id, exists := c.manifest[path]

	if !exists {
		return errs.InvalidOperation{
			Message: "asset not found: " + string(path),
		}
	}

	asset, err := c.store.Get(id)

	if err != nil {
		return errs.InvalidOperation{
			Message: "failed to retrieve asset from store: " + err.Error(),
		}
	}

	if err := c.store.Remove(id); err != nil {
		return errs.InvalidOperation{
			Message: "failed to remove asset from store: " + err.Error(),
		}
	}

	c.dealloc(asset)
	delete(c.manifest, path)
	return nil
}

// Get retrieves an asset from the cache by its AssetID. Returns an error if the AssetID is invalid or the asset is not found.
func (c *AssetCache[T]) Get(id AssetID) (T, error) {
	return c.store.Get(id)
}

// GetAssetID retrieves the AssetID associated with the given AssetPath. Returns an error if the asset is not found.
func (c *AssetCache[T]) GetAssetID(path AssetPath) (AssetID, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	id, exists := c.manifest[path]
	if !exists {
		return AssetID{}, errs.InvalidOperation{
			Message: "asset not found: " + string(path),
		}
	}

	return id, nil
}

// GetByPath retrieves an asset from the cache by its AssetPath. Returns an error if the asset is not found or the AssetID is invalid.
func (c *AssetCache[T]) GetByPath(path AssetPath) (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	id, exists := c.manifest[path]
	if !exists {
		var zero T
		return zero, errs.InvalidOperation{
			Message: "asset not found: " + string(path),
		}
	}

	return c.store.Get(id)
}

// NewAssetCache creates a new AssetCache with the specified allocation and deallocation functions and capacity for the underlying AssetStore.
func NewAssetCache[T any](allocate AllocateAssetFunc[T], deallocate DeallocateAssetFunc[T], capacity uint32) *AssetCache[T] {
	return &AssetCache[T]{
		store:    NewAssetStore[T](capacity),
		alloc:    allocate,
		dealloc:  deallocate,
		manifest: make(AssetManifest),
	}
}

// Assets is a high-level struct that manages different types of assets (e.g., images, sounds) using AssetCaches for each type.
type Assets struct {
	processors map[AssetType]AssetProcessor

	images *AssetCache[*ebiten.Image]
}

// NewAssets initializes a new Assets struct with AssetCaches for each asset type, setting up the appropriate allocation and deallocation functions.
func NewAssets() *Assets {
	return &Assets{
		processors: make(map[AssetType]AssetProcessor),
		images:     NewAssetCache[*ebiten.Image](images.AllocateImage, images.DeallocateImage, 1024),
	}
}

// Initialize registers AssetProcessors for each AssetType in the Assets struct. Returns an error if there are duplicate registrations or other issues.
func (a *Assets) Initialize() error {
	if err := registerProcessors(a.processors, a.images, "png", "jpg", "jpeg"); err != nil {
		return err
	}
	return nil
}

// Images returns the AssetCache for image assets.
func (a *Assets) Images() *AssetCache[*ebiten.Image] {
	return a.images
}

// Load loads assets from the given filesystem for the specified paths. Returns an error if any asset fails to load.
func (a *Assets) Load(filesystem fs.FS, paths ...AssetPath) error {
	for _, path := range paths {
		assetType := path.Type()

		if assetType == InvalidAssetType {
			return errs.InvalidArg{
				Message: "invalid asset type for path: " + string(path),
			}
		}

		processor, exists := a.processors[assetType]

		if !exists {
			return errs.InvalidArg{
				Message: "no processor registered for asset type: " + string(assetType),
			}
		}

		file, err := filesystem.Open(string(path))

		if err != nil {
			return errs.InvalidOperation{
				Message: "failed to open asset file: " + err.Error(),
			}
		}

		defer file.Close()

		if err := processor.allocate(path, file); err != nil {
			return errs.InvalidOperation{
				Message: "failed to allocate asset: " + err.Error(),
			}
		}
	}

	return nil
}

// Unload deallocates assets for the specified paths. Returns an error if any asset fails to unload.
func (a *Assets) Unload(paths ...AssetPath) error {
	for _, path := range paths {
		assetType := path.Type()

		if assetType == InvalidAssetType {
			return errs.InvalidArg{
				Message: "invalid asset type for path: " + string(path),
			}
		}

		processor, exists := a.processors[assetType]

		if !exists {
			return errs.InvalidArg{
				Message: "no processor registered for asset type: " + string(assetType),
			}
		}

		if err := processor.deallocate(path); err != nil {
			return errs.InvalidOperation{
				Message: "failed to deallocate asset: " + err.Error(),
			}
		}
	}

	return nil
}

// registerProcessors is a helper function that populates the processors map with AssetProcessors for the specified AssetTypes,
// checking for duplicates and returning an error if any are found.
func registerProcessors(table map[AssetType]AssetProcessor, loader AssetProcessor, types ...AssetType) error {
	for _, t := range types {
		if _, exists := table[t]; exists {
			return errs.DuplicateEntry{
				Message: "asset processor already registered for type: " + string(t),
			}
		}
		table[t] = loader
	}
	return nil
}
