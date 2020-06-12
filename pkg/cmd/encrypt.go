package cmd

import (
	"github.com/garyellis/secrets-ctl/pkg/secrets"
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

// Encrypt reads the input secret files, encrypts and rewrites them locally
func Encrypt(cmd *cobra.Command, args []string) error {
	secretsClient, err := secrets.New()
	if err != nil {
		return err
	}
	err = secretsClient.EncryptFilesWithVaultTransit(path, filter, backup)
	return err
}
