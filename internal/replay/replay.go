//go:build pcap

package replay

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/gopacket"
	gpcap "github.com/google/gopacket/pcap"
)

// Options configures the packet replay.
type Options struct {
	PcapFile  string
	Interface string
	Filter    string
	RateScale float64
	LoopCount int
}

// Run reads a PCAP file and replays the packets against the specified interface.
func Run(opts Options) error {
	handle, err := gpcap.OpenOffline(opts.PcapFile)
	if err != nil {
		return fmt.Errorf("failed to open PCAP file: %v", err)
	}
	defer handle.Close()

	if opts.Filter != "" {
		if err := handle.SetBPFFilter(opts.Filter); err != nil {
			return fmt.Errorf("failed to set BPF filter: %v", err)
		}
	}

	sendHandle, err := gpcap.OpenLive(opts.Interface, 65536, false, gpcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open interface %s for sending: %v", opts.Interface, err)
	}
	defer sendHandle.Close()

	fmt.Printf("Replaying [%s] on interface [%s] at %.1fx speed...\n", opts.PcapFile, opts.Interface, opts.RateScale)

	for loop := 0; loop < opts.LoopCount; loop++ {
		if opts.LoopCount > 1 {
			fmt.Printf("\n--- Loop %d / %d ---\n", loop+1, opts.LoopCount)
			// Reset the file handle for the next loop
			handle, _ = gpcap.OpenOffline(opts.PcapFile)
			if opts.Filter != "" {
				_ = handle.SetBPFFilter(opts.Filter)
			}
		}

		var prevTimestamp time.Time
		count := 0
		src := gopacket.NewPacketSource(handle, handle.LinkType())

		for packet := range src.Packets() {
			ts := packet.Metadata().Timestamp

			if !prevTimestamp.IsZero() && opts.RateScale > 0 {
				gap := ts.Sub(prevTimestamp)
				scaledGap := time.Duration(float64(gap) / opts.RateScale)
				if scaledGap > 0 && scaledGap < 5*time.Second {
					time.Sleep(scaledGap)
				}
			}
			prevTimestamp = ts

			if err := sendHandle.WritePacketData(packet.Data()); err != nil {
				fmt.Printf("Warning: failed to send packet #%d: %v\n", count+1, err)
			}
			count++
		}

		fmt.Printf("Loop %d complete. Replayed %d packets.\n", loop+1, count)
	}

	fmt.Println("\nReplay finished.")
	return nil
}

// ParseRateScale parses a rate multiplier string like "2x", "0.5x" into a float64.
func ParseRateScale(rate string) float64 {
	rate = strings.TrimSuffix(rate, "x")
	var f float64
	fmt.Sscanf(rate, "%f", &f)
	if f <= 0 {
		return 1.0
	}
	return f
}
