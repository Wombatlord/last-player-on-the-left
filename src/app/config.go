package app

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const DefaultConfig = "src/app/default_config.yaml"

type Subscription struct {
	Alias string `yaml:"alias"`
	Url   string `yaml:"url"`
}

type Config struct {
	Subs []Subscription `yaml:"subs"`
	Logs string         `yaml:"logs"`
}

func (c Config) GetByAlias(alias string) Subscription {
	for _, sub := range c.Subs {
		if sub.Alias == alias {
			return sub
		}
	}

	return Subscription{}
}

type ConfigFile struct {
	Path   string
	Config Config
}

var LoadedConfig Config

func (s *ConfigFile) Include(alias string, url string) error {
	s.Config.Subs = append(s.Config.Subs, Subscription{Alias: alias, Url: url})
	if err := s.Save(); err != nil {
		return err
	}
	return nil
}

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

func LoadConfig(path string) (*ConfigFile, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fileContent, err = os.ReadFile(DefaultConfig)
		if err != nil {
			log.Fatal(err)
		}
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
