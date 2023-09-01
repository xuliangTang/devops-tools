package cmd

import (
	"devops-tools/utils/helpers"
	"devops-tools/utils/remotes"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"syscall"
)

func init() {
	sshCmd.Flags().StringP("remote", "r", "", "find from remotes")
	sshCmd.Flags().StringP("server", "s", "", "set server host")
	sshCmd.Flags().StringP("username", "u", "", "set username")
	sshCmd.Flags().StringP("password", "p", "", "set password")
}

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "ssh connect",
	Run: func(cmd *cobra.Command, args []string) {
		var session *ssh.Session

		remoteName, err := cmd.Flags().GetString("remote")
		if err != nil {
			log.Fatalln(err)
		}

		if len(remoteName) > 0 { // 如果指定了remote配置则从配置创建连接
			getRemote, exist := remotes.GetRemote(remoteName)
			if !exist {
				log.Fatalln("remote not found:", remoteName)
			}
			session, err = helpers.SSHConnect(getRemote.Username, getRemote.Password, getRemote.Host, 22)
			if err != nil {
				log.Fatalln(err)
			}

		} else { // 没有指定配置或没有从配置中找到，则根据flag参数创建连接
			if session == nil {
				server := helpers.MustFlags(cmd, "server", "string").(string)
				user := helpers.MustFlags(cmd, "username", "string").(string)
				pwd, err := cmd.Flags().GetString("password")
				if err != nil {
					log.Fatalln(err)
				}

				if len(pwd) == 0 { // 判断没有指定-p 则进入隐藏密码输入
					connCount := 0
					for connCount < 3 {
						fmt.Println("entry password")
						getPwd, err := terminal.ReadPassword(syscall.Stdin)
						if err != nil {
							log.Fatalln(err)
						}

						session, err = helpers.SSHConnect(user, string(getPwd), server, 22)
						if err != nil {
							fmt.Println("error password, please try again")
							connCount++
							continue
						}
						break
					}
				} else {
					session, err = helpers.SSHConnect(user, pwd, server, 22)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}

		if session == nil {
			log.Fatalln("connect failed")
		}

		defer session.Close()
		session.Stdin = os.Stdin
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		var nodeShellModes = helpers.SSHTerminalModes()
		if err = session.RequestPty("", 0, 0, nodeShellModes); err != nil {
			log.Fatalln(err)
		}
		if err = session.Run("bash"); err != nil {
			log.Fatalln(err)
		}
	},
}
