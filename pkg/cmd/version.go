package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const cliName = "secrets-ctl"

var (
	GitCommit string
	Version   string
	BuildDate string
)

// VersionCmd prints the version
func VersionCmd() *cobra.Command {
	vCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cliName)
			fmt.Println("Version: ", Version)
			fmt.Println("GitCommit: ", GitCommit)
			fmt.Println("BuildDate: ", BuildDate)
		},
	}
	return vCmd
}
