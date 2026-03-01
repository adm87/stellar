package assets

import "strings"

// Asset represents a game asset path. The path is relative to the filesystem it is loaded from.
type Asset string

// String returns the string representation of the Asset, which is its path. This is useful for logging and debugging purposes.
func (a Asset) String() string {
	return string(a)
}

// Type returns the AssetType of the Asset by extracting the file extension from its path.
// For example, if the asset is "images/sprite.png", it will return "png". This is used to determine how to import the asset based on its type.
func (a Asset) Type() AssetType {
	ext := strings.LastIndex(string(a), ".")
	if ext == -1 || ext < strings.LastIndex(string(a), "/") || ext == len(string(a))-1 {
		return ""
	}
	return AssetType(string(a)[ext+1:])
}

// AssetType represents a type of asset, such as "jpeg", "png", "json", etc. This can be used to determine how to load the asset.
type AssetType string

// AssetImporter is a function type that defines how to import an asset. It takes an Asset and its raw byte data, and returns an error if the import fails.
//
// The AssetImporter is responsible for processing the raw data and converting it into a usable form for the game.
type AssetImporter func(asset Asset, rawData []byte) error
