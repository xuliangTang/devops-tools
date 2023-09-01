package cmd

import (
	"bytes"
	"devops-tools/utils/helpers"
	"devops-tools/utils/remotes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"path/filepath"
)

func init() {
	scpCmd.Flags().StringP("remote", "r", "", "set remote config name")
	scpCmd.Flags().StringP("source", "s", "", "set source path")
	scpCmd.Flags().StringP("dest", "d", "", "set remote path")
}

// 拷贝本地文件到远程主机
var scpCmd = &cobra.Command{
	Use:   "scp",
	Short: "copy local file to remote",
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := helpers.MustFlags(cmd, "remote", "string").(string)
		source := helpers.MustFlags(cmd, "source", "string").(string)
		dest := helpers.MustFlags(cmd, "dest", "string").(string)

		// 连接远程主机
		getRemote := remotes.MustGetRemote(remoteName)
		session, err := helpers.SSHConnect(getRemote.Username, getRemote.Password, getRemote.Host, 22)
		if err != nil {
			log.Fatalln(err)
		}
		defer session.Close()

		// 创建管道stdin
		in, err := session.StdinPipe()
		if err != nil {
			log.Fatalln(err)
		}

		// 读取本地文件
		b := helpers.MustLoadFile(source)

		// start启动远程scp -t代表接收模式
		if err = session.Start(fmt.Sprintf("scp -t %s", dest)); err != nil {
			log.Fatalln(err)
		}

		// 写入权限 大小 文件名
		if _, err = fmt.Fprintln(in, "C0644", int64(len(b)), filepath.Base(source)); err != nil {
			log.Fatalln(err)
		}

		// 向接收管道写入值
		n, err := io.Copy(in, bytes.NewReader(b))
		if err != nil {
			log.Fatalln(err)
		}

		// 写入结束符
		if _, err = fmt.Fprint(in, "\x00"); err != nil {
			log.Fatalln(err)
		}
		in.Close()

		if err = session.Wait(); err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("传输完成，共%d个字节\n", n)
	},
}
