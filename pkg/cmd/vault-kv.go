package cmd

import (
	vaultclient "github.com/garyellis/secrets-ctl/pkg/vault"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func VaultKVCmd() *cobra.Command {
	vaultCmd := &cobra.Command{
		Use:   "vault-kv",
		Short: "interacts with vault kv store",
	}
	vaultCmd.AddCommand(VaultKVWrite())
	return vaultCmd
}

func VaultKVWrite() *cobra.Command {
	vaultKVWriteCmd := &cobra.Command{
		Use:   "write",
		Short: "write secret files to vault kv secret engine",
		Long:  `decrypts encrypted secrets and writes them to vault kv secret engine`,
		RunE:  WriteVaultKV,
	}

	return vaultKVWriteCmd
}

func WriteVaultKV(cmd *cobra.Command, args []string) error {
	log.Infof("[cmd/vault-kv] reading secrets")
	c, err := DecryptFiles(path)
	if err != nil {
		return err
	}

	vault, err := vaultclient.NewClient()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("[cmd/vault-kv] writing secrets to vault kv")
	for _, secret := range c {
		log.Infof("[cmd/vault-kv] writing: %s", secret.Path)
		for _, i := range secret.Secret {
			log.Debugf("[cmd/vault-kv] uploading secret: %s", i.Data)
			err := vault.KV.Put(i.VaultKVPath, i.Data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
