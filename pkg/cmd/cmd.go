package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	path   string
	backup bool
	filter string
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secrets-ctl",
		Short:   "secrets utility cli",
		Version: "0.0.0",
	}

	cmd.PersistentFlags().StringVar(&path, "path", ".", "the secret yaml files path")
	cmd.PersistentFlags().StringVar(&filter, "filter", "secret.yaml", "the secret yaml files path")
	cmd.PersistentFlags().BoolVar(&backup, "backup", true, "backup the secrets file(s)")
	cmd.AddCommand(EncryptCmd())
	cmd.AddCommand(DecryptCmd())
	cmd.AddCommand(VaultKVCmd())
	return cmd
}

func init() {
	loglevelenv, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		loglevelenv = "info"
	}
	loglevel, err := logrus.ParseLevel(loglevelenv)
	if err != nil {
		loglevel = logrus.InfoLevel
	}
	logrus.SetLevel(loglevel)
}
