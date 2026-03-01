package assets

import (
	"io/fs"
	"sync"
)

// Loader is responsible for loading assets from the filesystem using registered importers.
// It ensures thread safety for the loader, not individual assets.
type Loader struct {
	filesystem fs.FS
	assets     []AssetPath
	mu         sync.Mutex
}

// NewLoader creates a new Loader with the given filesystem and assets to load.
func NewLoader(filesystem fs.FS, assets ...AssetPath) *Loader {
	return &Loader{
		filesystem: filesystem,
		assets:     assets,
	}
}

// Load loads all the assets using their respective importers.
func (l *Loader) Load() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, path := range l.assets {
		data, err := fs.ReadFile(l.filesystem, string(path))

		if err != nil {
			return FailedLoad{Asset: path, Err: err}
		}

		assetType := path.Type()

		if assetType == "" {
			return MissingType{Asset: path}
		}

		importer, exists := GetImporter(assetType)

		if !exists {
			return UnsupportedAssetType{AssetType: assetType}
		}

		if err := importer.Import(path, data); err != nil {
			return FailedImport{Asset: path, Err: err}
		}
	}

	return nil
}
