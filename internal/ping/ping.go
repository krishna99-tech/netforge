package ping

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

// Config holds ping configuration.
type Config struct {
	Count    int
	Interval time.Duration
	Timeout  time.Duration
}

// Run executes the ping operation.
func Run(host string, cfg Config) error {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return err
	}

	pinger.Count = cfg.Count
	pinger.Interval = cfg.Interval
	pinger.Timeout = cfg.Timeout
	pinger.SetPrivileged(true) // Required for Windows Administrator/Raw sockets

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%d\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	return pinger.Run()
}
