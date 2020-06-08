package cmd

import (
	"github.com/garyellis/secrets-ctl/pkg/config"
	"github.com/garyellis/secrets-ctl/pkg/util/fileutils"
	vaultclient "github.com/garyellis/secrets-ctl/pkg/vault"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// EncryptCmd reads the secret file and encrypts them
func EncryptCmd() *cobra.Command {
	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypts files on the given path",
		RunE:  Encrypt,
	}
	return encryptCmd
}

func Encrypt(cmd *cobra.Command, args []string) error {
	log.Infof("[cmd/encrypt] starting files encryption")
	c, err := EncryptFiles(path)
	if err != nil {
		return err
	}
	log.Infof("[cmd/encrypt] writing files")
	for _, secret := range c {
		log.Infof("[cmd/encrypt] writing: %s", secret.Path)
		err = fileutils.ToYamlFile(secret.Path, secret, backup)
		if err != nil {
			return err
		}
	}
	return nil
}

func EncryptFiles(path string) ([]*config.Config, error) {
	files, err := config.WalkMatch(path, filter)
	if err != nil {
		return nil, err
	}
	log.Info("[cmd/encrypt] procesing files: ", files)

	var encryptedConfigs []*config.Config
	for _, file := range files {
		secretConfig, err := EncryptFile(file)
		if err != nil {
			return nil, err
		}
		encryptedConfigs = append(encryptedConfigs, secretConfig)
	}
	return encryptedConfigs, nil
}

func EncryptFile(path string) (*config.Config, error) {
	log.Info("[cmd/encrypt] reading file: ", path)
	secretConfig, err := config.ReadConfig(path)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	vault, err := vaultclient.NewClient()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for i, secret := range secretConfig.Secret {
		for k, v := range secret.Data {
			log.Debugf("[cmd/encrypt] encrypting %s %s ", k, v)
			text, err := vault.Transit.Encrypt(
				secret.VaultMount,
				secret.VaultKey,
				[]byte(v.(string)),
			)
			if err != nil {
				log.Warnf("[cmd/encrypt] %s", err)
			}
			secretConfig.Secret[i].Data[k] = string(text)
			log.Debugf("[cmd/encrypt] encrypt %s %s", k, text)
		}
	}
	return secretConfig, nil
}
