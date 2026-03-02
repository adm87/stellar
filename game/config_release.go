//go:build release && !dev

package game

func NewConfig(version string, args Args) *Config {
	cfg := NewDefaultConfig(version, args)
	cfg.BuildMode = "release"
	return cfg
}
