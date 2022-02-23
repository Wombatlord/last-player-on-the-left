package app

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Path string
	Subs map[string]string
}

func (s *Config) Include(alias string, url string) error {
	s.Subs[alias] = url
	if err := s.Save(); err != nil {
		return err
	}
	return nil
}

func (s *Config) Save() error {
	content, err := yaml.Marshal(s.Subs)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.Path, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

var conf Config

func LoadConfig(path string) (*Config, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf.Subs = make(map[string]string)
	conf.Path = path
	err = yaml.Unmarshal(fileContent, conf.Subs)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}


