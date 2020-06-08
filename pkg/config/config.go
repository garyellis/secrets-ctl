package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Secret []Secret `yaml:"secret,omitempty" json:"secret,omitempty"`
	Path   string   `yaml:"-" json:"-"`
}

type Secret struct {
	VaultMount  string                 `yaml:"vault_mount,omitempty"`
	VaultKey    string                 `yaml:"vault_key,omitempty"`
	VaultKVPath string                 `yaml:"vault_kv_path,omitempty" json:"vault_kv_path,omitempty"`
	Data        map[string]interface{} `yaml:"data" json:"data,omitempty"`
}

func ReadConfig(path string) (*Config, error) {
	config := &Config{
		Path: path,
	}
	yamlfile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlfile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func ReadConfigFiles(path string) ([]*Config, error) {
	files, err := WalkMatch(path, "secret.yaml")
	log.Info("[config] found: ", files)
	if err != nil {
		return nil, err
	}
	config := []*Config{}
	for _, i := range files {
		c, err := ReadConfig(i)
		if err != nil {
			return nil, err
		}
		config = append(config, c)
	}
	return config, nil
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
