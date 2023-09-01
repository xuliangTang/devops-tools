package cmd

import (
	"devops-tools/utils/helpers"
	"devops-tools/utils/remotes"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	iptablesDropCmd.Flags().IntP("port", "p", 0, "set drop port")
}

// 禁用端口
var iptablesDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "drop filter port",
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := helpers.MustFlags(cmd, "remote", "string").(string)
		port := helpers.MustFlags(cmd, "port", "int").(int)

		// 连接远程主机
		getRemote := remotes.MustGetRemote(remoteName)
		session, err := helpers.SSHConnect(getRemote.Username, getRemote.Password, getRemote.Host, 22)
		if err != nil {
			log.Fatalln(err)
		}
		defer session.Close()

		execCmd := fmt.Sprintf("iptables -A INPUT -p tcp --dport %d -j DROP", port)
		if err = session.Run(execCmd); err != nil {
			log.Fatalln(err)
		}
	},
}
