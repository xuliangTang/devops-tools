package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "print the cpu and memory table",
	Run: func(cmd *cobra.Command, args []string) {
		// 设置表头
		var data [][]string
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"项目", "数量", "百分比"})

		// cpu信息
		cpuPercent, _ := cpu.Percent(time.Second, false)
		cpuCount, _ := cpu.Counts(true)
		data = append(data, []string{"CPU", fmt.Sprintf("%d核", cpuCount), fmt.Sprintf("%.2f%%", cpuPercent[0])})

		// 内存信息
		memInfo, _ := mem.VirtualMemory()
		data = append(data, []string{"内存", fmt.Sprintf("%dG", memInfo.Total/1024/1024/1024), fmt.Sprintf("%.2f%%", memInfo.UsedPercent)})

		// 渲染表格
		for _, v := range data {
			table.Append(v)
		}
		table.Render()
	},
}
