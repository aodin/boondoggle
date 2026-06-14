package boondoggle

import (
	"github.com/BurntSushi/toml"
)

// Config holds the settings used to build a site. It is typically loaded from a
// boondoggle.toml file, but the zero value (with DefaultConfig applied) is also
// usable.
type Config struct {
	Input     string `toml:"input"`     // Directory of markdown articles
	Output    string `toml:"output"`    // Directory to write the generated site
	Templates string `toml:"templates"` // Directory of HTML templates (optional)
	Previews  int    `toml:"previews"`  // Number of article previews on the index
	Feed      Feed   `toml:"feed"`      // RSS and Atom feed metadata
}

// DefaultConfig returns a Config populated with the default settings.
func DefaultConfig() Config {
	return Config{
		Input:    ".",
		Output:   "./dist",
		Previews: 4,
	}
}

// LoadConfig reads and decodes the TOML configuration file at path. Any setting
// absent from the file keeps its default value.
func LoadConfig(path string) (Config, error) {
	config := DefaultConfig()
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return config, err
	}
	return config, nil
}
