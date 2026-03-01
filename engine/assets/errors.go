package assets

// ------------------------------------------------------------------------------
// StoreFull Error
// ------------------------------------------------------------------------------

type StoreFull struct{}

func (e StoreFull) Error() string {
	return "cache is full"
}

// ------------------------------------------------------------------------------
// DuplicateAsset Error
// ------------------------------------------------------------------------------

type DuplicateAsset struct {
	Asset Asset
}

func (e DuplicateAsset) Error() string {
	return "duplicate asset: " + e.Asset.String()
}

// ------------------------------------------------------------------------------
// FailedImport Error
// ------------------------------------------------------------------------------

type FailedImport struct {
	Asset Asset
	Err   error
}

func (e FailedImport) Error() string {
	return "failed to import asset " + e.Asset.String() + ": " + e.Err.Error()
}

// ------------------------------------------------------------------------------
// FailedLoad Error
// ------------------------------------------------------------------------------

type FailedLoad struct {
	Asset Asset
	Err   error
}

func (e FailedLoad) Error() string {
	return "failed to load asset " + e.Asset.String() + ": " + e.Err.Error()
}

// ------------------------------------------------------------------------------
// MissingType Error
// ------------------------------------------------------------------------------

type MissingType struct {
	Asset Asset
}

func (e MissingType) Error() string {
	return "missing asset type for asset: " + e.Asset.String()
}

// ------------------------------------------------------------------------------
// UnsupportedAssetType Error
// ------------------------------------------------------------------------------

type UnsupportedAssetType struct {
	AssetType AssetType
}

func (e UnsupportedAssetType) Error() string {
	return "unsupported asset type: " + string(e.AssetType)
}
