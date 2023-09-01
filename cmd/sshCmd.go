package cmd

import (
	"devops-tools/utils/helpers"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net"
	"os"
	"syscall"
)

func init() {
	sshCmd.Flags().StringP("server", "s", "", "set server host")
	sshCmd.Flags().StringP("username", "u", "", "set username")
	sshCmd.Flags().StringP("password", "p", "", "set password")
}

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "ssh connect",
	Run: func(cmd *cobra.Command, args []string) {
		server := helpers.MustFlags(cmd, "server", "string").(string)
		user := helpers.MustFlags(cmd, "username", "string").(string)
		pwd, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalln(err)
		}

		var session *ssh.Session
		if len(pwd) == 0 { // 判断没有指定-p 则进入隐藏密码输入
			connCount := 0
			for connCount < 3 {
				fmt.Println("entry password")
				getPwd, err := terminal.ReadPassword(syscall.Stdin)
				if err != nil {
					log.Fatalln(err)
				}

				session, err = sshConnect(user, string(getPwd), server, 22)
				if err != nil {
					fmt.Println("error password, please try again")
					connCount++
					continue
				}
				break
			}
		} else {
			session, err = sshConnect(user, pwd, server, 22)
			if err != nil {
				log.Fatalln(err)
			}
		}

		if session == nil {
			log.Fatalln("connect failed")
		}

		defer session.Close()
		session.Stdin = os.Stdin
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		var nodeShellModes = ssh.TerminalModes{
			ssh.ECHO:          1,     // enable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}
		if err = session.RequestPty("", 0, 0, nodeShellModes); err != nil {
			log.Fatalln(err)
		}
		if err = session.Run("bash"); err != nil {
			log.Fatalln(err)
		}
	},
}

// SSHConnect 获取SSH连接
func sshConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallback,
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}
