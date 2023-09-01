package cmd

import (
	"devops-tools/utils/helpers"
	"devops-tools/utils/remotes"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	shellCmd.Flags().StringP("remote", "r", "", "find from remotes")
	shellCmd.Flags().StringP("command", "c", "", "set execute command")
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "execute shell command",
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := helpers.MustFlags(cmd, "remote", "string").(string)
		command := helpers.MustFlags(cmd, "command", "string").(string)

		getRemote, exist := remotes.GetRemote(remoteName)
		if !exist {
			log.Fatalln("no such remote:", remoteName)
		}

		session, err := helpers.SSHConnect(getRemote.Username, getRemote.Password, getRemote.Host, 22)
		if err != nil {
			log.Fatalln(err)
		}
		defer session.Close()
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		if err = session.Run(command); err != nil {
			log.Fatalln(err)
		}
	},
}
