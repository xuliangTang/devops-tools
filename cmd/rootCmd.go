package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{}

func RunCmd() {
	rootCmd.AddCommand(versionCmd, cpuCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
