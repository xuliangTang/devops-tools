package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Version = "v0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the devops-tool version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
