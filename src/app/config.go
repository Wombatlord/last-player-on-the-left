package app

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

// Subscription represents a single alias <-> url pair. These are the items that show up
// in the feeds menu
type Subscription struct {
	Alias string `yaml:"alias"`
	Url   string `yaml:"url"`
}

// Config represents all the configuration contained in a config file. It
// specifies the config schema.
type Config struct {
	Subs []Subscription `yaml:"subs"`
	Logs string         `yaml:"logs"`
}

// GetByAlias returns the Subscription associated to the passed alias
func (c Config) GetByAlias(alias string) Subscription {
	for _, sub := range c.Subs {
		if sub.Alias == alias {
			return sub
		}
	}

	return Subscription{}
}

// ConfigFile represents a file containing configuration
type ConfigFile struct {
	Path   string
	Config Config
}

var LoadedConfig Config

// DefaultConfig is the configuration that will be saved in the configured
// path if no file is found on LoadConfig
var DefaultConfig = ConfigFile{
	Path: GetPath(),
	Config: Config{
		Subs: []Subscription{},
		Logs: "logs/log.txt",
	},
}

// Include allows the application to add subscriptions to the config
func (s *ConfigFile) Include(alias string, url string) error {
	s.Config.Subs = append(s.Config.Subs, Subscription{Alias: alias, Url: url})
	if err := s.Save(); err != nil {
		return err
	}
	return nil
}

// Save updates the config file that the config was loaded from with any changes
func (s *ConfigFile) Save() error {
	content, err := yaml.Marshal(s.Config)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.Path, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

var conf ConfigFile
var confVals Config

// LoadConfig checks the environment for a LAST_CONFIG_PATH_ON_THE_LEFT variable for a
// user supplied config path, if empty it puts it into $HOME/.config/LastPlayer. If no
// config is found at the specified path, the default config is placed there.
func LoadConfig() (*ConfigFile, error) {
	path := GetPath()
	fileContent, err := os.ReadFile(path)
	if err != nil {
		LoadedConfig = DefaultConfig.Config
		err = DefaultConfig.Save()
		if err != nil {
			log.Fatal(err)
		}
		return &DefaultConfig, nil
	}
	conf.Path = path
	err = yaml.Unmarshal(fileContent, &confVals)
	if err != nil {
		return nil, err
	}
	conf.Config = confVals
	LoadedConfig = confVals

	return &conf, nil
}

// GetPath returns the config path as specified in env variable LAST_CONFIG_PATH_ON_THE_LEFT
// If empty the fallback is ~/.config/LastPlayer/config.yaml
func GetPath() string {
	path := os.Getenv("LAST_CONFIG_PATH_ON_THE_LEFT")
	if path == "" {
		confDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal(err)
		}
		return confDir + filepath.FromSlash("/last_player/config.yaml")
	}
	return path
}
