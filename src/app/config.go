package app

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const DefaultConfig = "src/app/default_config.yaml"

type Config struct {
	Subs map[string]string `yaml:"subs"`
	Logs string            `yaml:"logs"`
}

type ConfigFile struct {
	Path   string
	Config Config
}

func (s *ConfigFile) Include(alias string, url string) error {
	if s.Config.Subs == nil {
		s.Config.Subs = make(map[string]string)
	}
	s.Config.Subs[alias] = url
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
	conf.Config = confVals
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
