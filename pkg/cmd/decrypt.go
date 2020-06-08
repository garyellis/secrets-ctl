package cmd

import (
	"github.com/garyellis/secrets-ctl/pkg/config"
	"github.com/garyellis/secrets-ctl/pkg/util/fileutils"
	vaultclient "github.com/garyellis/secrets-ctl/pkg/vault"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DecryptCmd reads the secret file and encrypts them
func DecryptCmd() *cobra.Command {
	encryptCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypts files on the given path",
		RunE:  Decrypt,
	}
	return encryptCmd
}

// Decrypt reads the input secret files, decrypts and rewrites them locally
func Decrypt(cmd *cobra.Command, args []string) error {
	log.Infof("[cmd/decrypt] starting files decryption")
	c, err := DecryptFiles(path)
	if err != nil {
		return err
	}

	log.Infof("[cmd/decrypt] writing files")
	for _, secret := range c {
		log.Infof("[cmd/decrypt] writing: %s", secret.Path)
		err = fileutils.ToYamlFile(secret.Path, secret, backup)
		if err != nil {
			return err
		}
	}
	return nil
}

// DecryptFiles reads the input list of files, and decrypts them into a slice of secrets config
func DecryptFiles(path string) ([]*config.Config, error) {
	files, err := config.WalkMatch(path, filter)
	if err != nil {
		return nil, err
	}
	log.Info("[cmd/decrypt] procesing files: ", files)
	var decryptedConfigs []*config.Config

	for _, file := range files {
		secretConfig, err := DecryptFile(file)
		if err != nil {
			return nil, err
		}
		decryptedConfigs = append(decryptedConfigs, secretConfig)
	}
	return decryptedConfigs, nil
}

// DecryptFile reads the the input file and decrypts the encrypted values into a secret config
func DecryptFile(path string) (*config.Config, error) {
	log.Info("[cmd/decrypt] reading file: ", path)
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
			log.Debugf("[cmd/decrypt] decrypting %s %s ", k, v)
			text, err := vault.Transit.Decrypt(
				secret.VaultMount,
				secret.VaultKey,
				[]byte(v.(string)),
			)
			if err != nil {
				log.Warn(err)
			}
			secretConfig.Secret[i].Data[k] = string(text)
			log.Debugf("[cmd/decrypt] decrypted %s %s", k, text)
		}
	}
	return secretConfig, nil
}
