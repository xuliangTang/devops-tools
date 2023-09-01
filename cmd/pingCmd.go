package cmd

import (
	"devops-tools/utils/helpers"
	"devops-tools/utils/remotes"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"runtime"
)

func init() {
	pingCmd.Flags().StringP("remote", "r", "", "set remote name")
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping remote host",
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := helpers.MustFlags(cmd, "remote", "string").(string)
		getRemote := remotes.MustGetRemote(remoteName)

		pinger, err := ping.NewPinger(getRemote.Host)
		if err != nil {
			panic(err)
		}

		// Listen for Ctrl-C.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for _ = range c {
				pinger.Stop()
			}
		}()

		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}

		pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
		}

		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}

		if runtime.GOOS == "windows" {
			pinger.SetPrivileged(true)
		}

		pinger.Count = 5 // 设置ping次数

		fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
		err = pinger.Run()
		if err != nil {
			panic(err)
		}
	},
}
