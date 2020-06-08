package fileutils

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ToYamlFile writes the content to a file
func ToYamlFile(path string, content interface{}, backup bool) error {
	yamlContent, err := yaml.Marshal(&content)
	if err != nil {
		return err
	}

	// backup the file before writing it
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if backup {
			log.Infof("[util/fileutils] backup file %s to %s_backup", path, path)
			backupContent, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(fmt.Sprintf("%s_backup", path), backupContent, 0644)
			if err != nil {
				return err
			}
		}
	}
	err = ioutil.WriteFile(path, yamlContent, 0644)
	return err
}
