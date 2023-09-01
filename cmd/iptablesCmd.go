package cmd

import (
	"devops-tools/utils/helpers"
	"devops-tools/utils/remotes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func init() {
	iptablesCmd.PersistentFlags().StringP("remote", "r", "", "set remote config name")
	iptablesCmd.PersistentFlags().StringP("table", "t", "", "set table")
	// 添加子命令
	iptablesCmd.AddCommand(iptablesDropCmd)
}

var iptablesCmd = &cobra.Command{
	Use:   "iptables",
	Short: "iptables show",
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := helpers.MustFlags(cmd, "remote", "string").(string)
		table := helpers.MustFlags(cmd, "table", "string").(string)

		// 连接远程主机
		getRemote := remotes.MustGetRemote(remoteName)
		session, err := helpers.SSHConnect(getRemote.Username, getRemote.Password, getRemote.Host, 22)
		if err != nil {
			log.Fatalln(err)
		}
		defer session.Close()

		out, _ := session.StdoutPipe()
		execCmd := fmt.Sprintf("iptables -t %s -nvL --line-number", table)
		if err = session.Run(execCmd); err != nil {
			log.Fatalln(err)
		}

		b, err := io.ReadAll(out)
		if err != nil {
			log.Fatalln(err)
		}
		//render(string(b))
		fmt.Println(string(b))
	},
}

var chanInfo = regexp.MustCompile(`^Chain\s*(INPUT|FORWARD|OUTPUT)`)
var headerInfo = regexp.MustCompile(`^num\s+pkts\s+bytes`)

// 以表格方式输出
func render(str string) {
	list := strings.Split(str, "\n")
	var header []string
	var data [][]string
	begin := true
	for _, item := range list {
		if item = strings.Trim(item, " "); item != "" {
			if isChanInfo(item) {
				begin = true
			} else {
				begin = false
			}
			if begin { // 需要打印
				if len(header) > 0 {
					printtable(header, data)
					header = []string{}
					data = [][]string{}
				}
				fmt.Println("链信息", item)
			}
			if isHeaderInfo(item) {
				header = append(header, filterlist(item)...)
			}
			if !isChanInfo(item) && !isHeaderInfo(item) {
				data = append(data, filterlist(item))
			}

		}
	}
	printtable(header, data)
}

func isChanInfo(str string) bool {
	return chanInfo.MatchString(str)
}

func isHeaderInfo(str string) bool {
	return headerInfo.MatchString(str)
}

func printtable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func filterlist(str string) []string {
	list := strings.Split(str, " ")
	var ret []string
	for _, item := range list {
		if item = strings.Trim(item, " "); item != "" {
			ret = append(ret, item)
		}
	}
	return ret
}
