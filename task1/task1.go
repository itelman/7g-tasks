package task1

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

/*
Network traffic analytics

Steps:
1. Read packets from given interface (example: enp3s0)
2. Apply a kernel-level BPF filter for a given host IP: only accept packets involving that IP
(example: 192.168.0.24)
3. Maintain sliding window of last N packets (example: 100), keep a frequency distribution of
source IPs in that window
4. Compute Shannon entropy over that window every T seconds to measure traffic diversity (example: 3)

This is used in:
- DDoS detection
- Port scan detection
- Traffic anomaly detection
*/

/*
Interview Talking Points, Advanced Upgrade:

- Instead of IP entropy, compute entropy over:
- - Destination ports
- - Packet sizes
- - 5-tuples
*/

func Run() {
	if len(os.Args) < 5 {
		log.Fatal("Usage: ./entropy <interface> <filterIP> <windowSize> <periodSeconds>")
	}

	iface := os.Args[1]
	filterIP := os.Args[2]
	windowSize, _ := strconv.Atoi(os.Args[3])
	periodSec, _ := strconv.Atoi(os.Args[4])

	// Open interface
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Kernel-level BPF filtering (better performance):
	// Instead of the code filtering packets manually, the Linux kernel drops unwanted packets.
	// Advantages:
	// - Reduces CPU usage
	// - Avoids unnecessary memory copies
	// - Improves performance
	err = handle.SetBPFFilter("host " + filterIP)
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Sliding window mechanism

	// 'window' stores the last N source IPs (FIFO queue)
	window := make([]string, 0, windowSize)
	// 'counts' stores how many times each IP appears in the window
	counts := make(map[string]int)

	// Packet handling and entropy calculation are executed in parallel. This ensures:
	// - Packet processing is continuous
	// - Printing results doesn't block capture

	ticker := time.NewTicker(time.Duration(periodSec) * time.Second)
	go func() {
		for range ticker.C {
			h := calculateEntropy(counts, len(window))
			fmt.Printf("Entropy: %.4f (window size: %d)\n", h, len(window))
		}
	}()

	// Packet handling (O(N*1) complexity):
	// N - number of packets
	for packet := range packetSource.Packets() {
		// Extract IPv4 layer
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}

		ip := ipLayer.(*layers.IPv4)

		// Get source IP
		src := ip.SrcIP.String()

		// Append source IP to window
		window = append(window, src)
		// Increment count in map
		counts[src]++

		// Maintain sliding window
		// If window exceeds size N:
		if len(window) > windowSize {
			// The oldest packet is removed
			old := window[0]
			window = window[1:]

			counts[old]--
			// If count becomes zero:
			if counts[old] == 0 {
				delete(counts, old)
			}
		}

		// At any moment 'counts' map contains exact frequency distribution of current sliding window
		// This makes entropy calculation correct and efficient
	}
}

// Time complexity - O(K)
// K - number of unique IPs
// Worst case: O(windowSize) - if K equals window size

// Higher entropy - more diverse traffic
// Lower entropy - more uniform traffic
func calculateEntropy(counts map[string]int, total int) float64 {
	if total == 0 {
		return 0
	}

	var entropy float64

	// Computation of Shannon entropy over packet source IP distribution:

	// 'count' = number of packets from one IP; count of value x in window
	// 'total' = window size
	// 'p' = probability of that IP
	for _, count := range counts {
		p := float64(count) / float64(total)
		entropy -= p * math.Log2(p)
	}

	return entropy
}
