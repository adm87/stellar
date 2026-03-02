//go:build dev || !release

package game

func NewConfig(version string, args Args) *Config {
	cfg := NewDefaultConfig(version, args)
	cfg.BuildMode = "dev"
	return cfg
}
