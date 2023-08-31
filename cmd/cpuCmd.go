package cmd

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/spf13/cobra"
	"time"
)

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "print the cpu percent",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			p, _ := cpu.Percent(time.Second, false)
			fmt.Printf("\r%.2f%%", p[0])
			time.Sleep(time.Second * 1)
		}
	},
}
