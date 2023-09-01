package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{}

func RunCmd() {
	rootCmd.AddCommand(versionCmd, cpuCmd, infoCmd, sshCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
