package secrets

import (
	"io/ioutil"
	"strings"

	"github.com/garyellis/secrets-ctl/pkg/util/fileutils"
	vaultclient "github.com/garyellis/secrets-ctl/pkg/vault"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type SecretsClient struct {
	Secrets     []*Secret
	VaultClient *vaultclient.Client
}

type Secret struct {
	SecretConfig SecretConfig `yaml:"secretconfig,omitempty" json:"secretconfig,omitempty"`
}

type SecretConfig struct {
	Encryption Encryption   `yaml:"encryption,omitempty" json:"encryption,omitempty"`
	Secrets    []SecretData `yaml:"secrets,omitempty" json:"secrets,omitempty"`
	Path       string       `yaml:"-" json:"-"`
}

type Encryption struct {
	VaultTransitMount string `yaml:"transit_mount,omitempty"`
	VaultTransitKey   string `yaml:"transit_key,omitempty"`
}

type SecretData struct {
	VaultKVPath string                 `yaml:"vault_kv_path,omitempty" json:"vault_kv_path,omitempty"`
	Data        map[string]interface{} `yaml:"data" json:"data,omitempty"`
}

func New() (*SecretsClient, error) {
	vaultClient, err := vaultclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &SecretsClient{
		VaultClient: vaultClient,
	}, nil
}

func (s *SecretsClient) EncryptFilesWithVaultTransit(folder, filter string, backup bool) error {
	// read the yaml configuration files into a slice of secret config files
	err := s.ReadSecretConfigFolder(folder, filter)
	if err != nil {
		return err
	}
	err = s.VaultTransitEncrypt()
	if err != nil {
		return err
	}
	for i := range s.Secrets {
		err = fileutils.ToYamlFile(s.Secrets[i].SecretConfig.Path, s.Secrets[i], backup)
		if err != nil {
			log.Error("[secrets/secrets] ", err)
		}
	}
	return nil
}

func (s *SecretsClient) DecryptFilesWithVaultTransit(folder, filter string, backup bool) error {
	// read the yaml configuration files into a slice of secret config files
	err := s.ReadSecretConfigFolder(folder, filter)
	if err != nil {
		return err
	}
	err = s.VaultTransitDecrypt()
	if err != nil {
		return err
	}
	for i := range s.Secrets {
		err = fileutils.ToYamlFile(s.Secrets[i].SecretConfig.Path, s.Secrets[i], backup)
		if err != nil {
			log.Error("[secrets/secrets] ", err)
		}
	}
	return nil
}

func (s *SecretsClient) WriteFilesToVaultKV(folder, filter string) error {
	// read the yaml configuration files into a slice of secret config files
	err := s.ReadSecretConfigFolder(folder, filter)
	if err != nil {
		return err
	}
	err = s.VaultTransitDecrypt()
	if err != nil {
		return err
	}

	// write the secrets to vault kv
	for _, secretConfig := range s.Secrets {
		for _, secret := range secretConfig.SecretConfig.Secrets {
			err = s.VaultClient.KV.Put(
				secret.VaultKVPath,
				secret.Data,
			)
			if err != nil {
				log.Error("[secrets/secrets]", err)
				return err
			}
		}
	}
	return nil
}

func (s *SecretsClient) VaultExportKV(key, transitmount, transitkey, filename string, encrypt bool) error {
	s.Secrets = append(s.Secrets, &Secret{
		SecretConfig: SecretConfig{
			Encryption: Encryption{
				VaultTransitKey:   transitkey,
				VaultTransitMount: transitmount,
			},
			Path: filename,
		},
	})

	// get the kv keys
	log.Debugf("[secrets/secrets] fetching keys for %s", key)
	keys, err := s.VaultClient.GetKeys(key)
	if err != nil {
		return err
	}
	log.Debugf("[secrets/secrets] got keys: %s", strings.Join(keys, ","))
	for _, i := range keys {
		secret, err := s.VaultClient.KV.ReadSecretData(i)
		if err != nil {
			return err
		}

		data := secret["data"].(map[string]interface{})
		s.Secrets[0].SecretConfig.Secrets = append(s.Secrets[0].SecretConfig.Secrets, SecretData{
			VaultKVPath: i,
			Data:        data,
		})
	}

	// write the file
	err = fileutils.ToYamlFile(s.Secrets[0].SecretConfig.Path, s.Secrets[0], false)
	return err
}

func (s *SecretsClient) ReadSecretConfigFile(file string) error {
	secretConfig := &Secret{}

	yamlfile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error(err)
		return err
	}
	err = yaml.Unmarshal(yamlfile, secretConfig)
	if err != nil {
		log.Error(err)
	}
	secretConfig.SecretConfig.Path = file
	s.Secrets = append(s.Secrets, secretConfig)
	return err
}

func (s *SecretsClient) ReadSecretConfigFolder(folder string, filter string) error {
	files, err := fileutils.WalkMatch(folder, filter)
	if err != nil {
		return err
	}
	log.Info("[secrets/secrets] procesing files: ", files)
	for _, file := range files {
		err = s.ReadSecretConfigFile(file)
		if err != nil {
			log.Error("[secrets/secrets]", err)
			return err
		}
	}
	return nil
}

func (s *SecretsClient) VaultTransitEncrypt() error {
	for i, secretConfig := range s.Secrets {
		for i2, secret := range secretConfig.SecretConfig.Secrets {
			for k, v := range secret.Data {
				ciphertext, err := s.VaultClient.Transit.Encrypt(
					secretConfig.SecretConfig.Encryption.VaultTransitMount,
					secretConfig.SecretConfig.Encryption.VaultTransitKey,
					[]byte(v.(string)),
				)
				if err != nil {
					log.Warnf("[secrets/secrets] %s", err)
				}
				s.Secrets[i].SecretConfig.Secrets[i2].Data[k] = string(ciphertext)
			}
		}
	}
	return nil
}

func (s *SecretsClient) VaultTransitDecrypt() error {
	for i, secretConfig := range s.Secrets {
		for i2, secret := range secretConfig.SecretConfig.Secrets {
			for k, v := range secret.Data {
				plaintext, err := s.VaultClient.Transit.Decrypt(
					secretConfig.SecretConfig.Encryption.VaultTransitMount,
					secretConfig.SecretConfig.Encryption.VaultTransitKey,
					[]byte(v.(string)),
				)
				if err != nil {
					log.Warnf("[secrets/secrets] %s", err)
				}
				s.Secrets[i].SecretConfig.Secrets[i2].Data[k] = string(plaintext)
			}
		}
	}
	return nil
}
