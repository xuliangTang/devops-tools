package cmd

import (
	"devops-tools/utils/helpers"
	"devops-tools/utils/remotes"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"regexp"
	"strconv"
)

func init() {
	iptablesDeleteCmd.Flags().StringP("chain", "c", "INPUT", "set chain")
	iptablesDeleteCmd.Flags().StringP("line", "l", "", "set line number")
}

// 根据行号删除规则
var iptablesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete by line number",
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := helpers.MustFlags(cmd, "remote", "string").(string)
		chain := helpers.MustFlags(cmd, "chain", "string").(string)
		line := helpers.MustFlags(cmd, "line", "string").(string)

		// 连接远程主机
		getRemote := remotes.MustGetRemote(remoteName)
		session, err := helpers.SSHConnect(getRemote.Username, getRemote.Password, getRemote.Host, 22)
		if err != nil {
			log.Fatalln(err)
		}
		defer session.Close()
		in, _ := session.StdinPipe()
		if err = session.Shell(); err != nil {
			log.Fatalln(err)
		}

		// 解析行号范围，批量删除
		lineRange := parseRange(line)
		for i := lineRange.start; i <= lineRange.end; i++ {
			execCmd := fmt.Sprintf("iptables -D %s %d", chain, lineRange.start)
			in.Write([]byte(execCmd + "\n"))
		}
	},
}

type lineRange struct {
	start int
	end   int
}

// 解析行号范围字符串，如2-5
func parseRange(l string) *lineRange {
	ret := &lineRange{}
	reg := regexp.MustCompile("^(?P<start>\\d*)-?(?P<end>\\d*)$")
	if match := reg.FindStringSubmatch(l); len(match) > 0 {
		for i, name := range reg.SubexpNames() {
			if i != 0 && match[i] != "" {
				num, _ := strconv.Atoi(match[i])
				if name == "start" {
					ret.start = num
				}
				if name == "end" {
					ret.end = num
				}

			}
		}
	}
	if ret.start < 0 {
		ret.start = 0
	}
	if ret.end < ret.start {
		ret.end = ret.start
	}
	return ret
}
