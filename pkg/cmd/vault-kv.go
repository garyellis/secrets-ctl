package cmd

import (
	"github.com/garyellis/secrets-ctl/pkg/secrets"
	vaultclient "github.com/garyellis/secrets-ctl/pkg/vault"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	kvSearchPath    string
	kvMount         string
	transitMount    string
	transitKey      string
	outputFilename  string
	encryptExported bool
)

func VaultKVCmd() *cobra.Command {
	vaultCmd := &cobra.Command{
		Use:   "vault-kv",
		Short: "interact with vault kv secret engine",
	}
	vaultCmd.PersistentFlags().StringVar(&kvSearchPath, "path", "/secret/data/", "the kv search path")
	vaultCmd.AddCommand(VaultKVWriteCmd())
	vaultCmd.AddCommand(VaultKVListCmd())
	vaultCmd.AddCommand(VaultKVExportCmd())
	return vaultCmd
}

func VaultKVListCmd() *cobra.Command {
	vaultKVListCmd := &cobra.Command{
		Use:   "list",
		Short: "lists vault kv keys",
		RunE:  VaultKVList,
	}
	return vaultKVListCmd
}

func VaultKVExportCmd() *cobra.Command {
	vaultKVExportCmd := &cobra.Command{
		Use:   "export",
		Short: "export kv secrets to a yaml config file",
		RunE:  VaultKVExport,
	}
	vaultKVExportCmd.PersistentFlags().StringVar(&transitMount, "transit-mount", "/transit", "the transit engine mount path")
	vaultKVExportCmd.PersistentFlags().StringVar(&transitKey, "transit-key", "${TRANSIT_KEY}", "transit engine encryption key")
	vaultKVExportCmd.PersistentFlags().StringVar(&outputFilename, "out", "vault-kv-export.yaml", "the output filename")
	vaultKVExportCmd.PersistentFlags().BoolVar(&encryptExported, "encrypt", false, "encrypt secrets data")
	return vaultKVExportCmd
}

func VaultKVWriteCmd() *cobra.Command {
	vaultKVWriteCmd := &cobra.Command{
		Use:   "write",
		Short: "write secret yaml files to vault",
		Long:  `writes yaml secret config files to vault kv secret engine`,
		RunE:  VaultKVWrite,
	}

	return vaultKVWriteCmd
}

func VaultKVExport(cmd *cobra.Command, args []string) error {
	log.Infof("[cmd/vault-kv] exporting vault kv secrets")
	secretsClient, err := secrets.New()
	if err != nil {
		return err
	}
	err = secretsClient.VaultExportKV(kvSearchPath, transitMount, transitKey, outputFilename, encryptExported)
	return err
}

func VaultKVWrite(cmd *cobra.Command, args []string) error {
	log.Infof("[cmd/vault-kv] writing secret config files to vault kv")
	secretsClient, err := secrets.New()
	if err != nil {
		return err
	}
	err = secretsClient.WriteFilesToVaultKV(path, filter)
	return err
}

func VaultKVList(cmd *cobra.Command, args []string) error {
	log.Infof("[cmd/vault-kv] listing vault kv secrets")
	vault, err := vaultclient.NewClient()
	if err != nil {
		log.Error(err)
		return err
	}
	keys, err := vault.KV.GetKeys(kvSearchPath)
	if err != nil {
		return err
	}
	for _, key := range keys {
		log.Infof("[cmd/vault-kv]  key: %s", key)
	}
	return nil
}
