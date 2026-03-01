package assets

import "sync"

var (
	importers  = make(map[AssetType]AssetImporter)
	importerMu sync.Mutex
)

// RegisterImporter registers an AssetImporter for a specific AssetType.
// If an importer for the given AssetType already exists, it panics to prevent overwriting the existing importer.
func RegisterImporter(assetType AssetType, importer AssetImporter) {
	importerMu.Lock()
	defer importerMu.Unlock()

	if _, exists := importers[assetType]; exists {
		panic("importer for asset type " + string(assetType) + " already registered")
	}
	importers[assetType] = importer
}

// GetImporter retrieves the AssetImporter for a given AssetType. It returns the importer and a boolean indicating whether the importer exists.
func GetImporter(assetType AssetType) (AssetImporter, bool) {
	importerMu.Lock()
	defer importerMu.Unlock()

	importer, exists := importers[assetType]
	return importer, exists
}
