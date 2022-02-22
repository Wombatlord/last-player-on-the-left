package app

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"os"
)

type Subscription struct {
	Url   string
	Alias string
}

type SubscriptionsConf struct {
	file *os.File
	Subs []Subscription
}

func (s *SubscriptionsConf) Include(alias string, url string) {
	s.Subs = append(s.Subs, Subscription{Alias: alias, Url: url})
}

func (s *SubscriptionsConf) Save() error {
	content, err := yaml.Marshal(s.Subs)
	if err != nil {
		return err
	}
	_, err = s.file.Write(content)
	if err != nil {
		return err
	}
	return nil
}

var subs SubscriptionsConf

func LoadSubs(file *os.File) (*SubscriptionsConf, error) {
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	fileContent := scanner.Bytes()
	err := yaml.Unmarshal(fileContent, subs.Subs)
	if err != nil {
		return nil, err
	}

	return &subs, nil
}
