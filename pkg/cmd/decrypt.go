package cmd

import (
	"github.com/garyellis/secrets-ctl/pkg/secrets"
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
	secretsClient, err := secrets.New()
	if err != nil {
		return err
	}
	err = secretsClient.DecryptFilesWithVaultTransit(path, filter, backup)
	return err
}
