package task2

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"time"
)

/*
Objective:

Receives info from the sniffer, and keeps stats: how many packets and total bytes per IP
Every 5 seconds, prints: which IP had the most packets, and which transferred the most data

Steps:

1. Receives JSON metadata from sniffer.
2. Updates stats map.
3. Every 5 seconds → prints top IPs.
*/

type Stats struct {
	IP      string
	Packets int
	Bytes   int
}

// Each packet updates stats
// Time complexity: O(1)
func update(stats map[string]*Stats, ip string, length int) {
	if _, ok := stats[ip]; !ok {
		stats[ip] = &Stats{IP: ip}
	}
	stats[ip].Packets++
	stats[ip].Bytes += length
}

func printStats(stats map[string]*Stats) {
	hosts := make([]Stats, 0, len(stats))

	for _, v := range stats {
		hosts = append(hosts, *v)
	}

	fmt.Printf("Unique IPs observed: %d\n\n", len(hosts))

	// Sort by packets
	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i].Packets > hosts[j].Packets
	})

	fmt.Println("Top 5 by packets:")
	for i := 0; i < 5 && i < len(hosts); i++ {
		fmt.Printf("%d. %s (%d packets)\n",
			i+1,
			hosts[i].IP,
			hosts[i].Packets,
		)
	}

	// Sort by bytes
	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i].Bytes > hosts[j].Bytes
	})

	fmt.Println("\nTop 5 by bytes:")
	for i := 0; i < 5 && i < len(hosts); i++ {
		fmt.Printf("%d. %s (%d bytes)\n",
			i+1,
			hosts[i].IP,
			hosts[i].Bytes,
		)
	}
}

func RunAnalyzer() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./analyzer <listen_host:port>")
	}

	listenAddr := os.Args[1]
	stats := make(map[string]*Stats)

	// Listen on UDP
	addr, _ := net.ResolveUDPAddr("udp", listenAddr)
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	// Every 5 seconds:
	// - Iterate through map
	// - Find:
	// - - IP with maximum packets
	// - - IP with maximum bytes

	// Time complexity: O(N)
	// N = number of unique IPs.
	go func() {
		for {
			time.Sleep(5 * time.Second)

			var maxPacketsIP string
			var maxBytesIP string
			maxPackets := 0
			maxBytes := 0

			for ip, s := range stats {
				if s.Packets > maxPackets {
					maxPackets = s.Packets
					maxPacketsIP = ip
				}
				if s.Bytes > maxBytes {
					maxBytes = s.Bytes
					maxBytesIP = ip
				}
			}

			log.Println("===== Statistics =====")
			fmt.Printf("Top by packets: %s (%d packets)\n", maxPacketsIP, maxPackets)
			fmt.Printf("Top by bytes: %s (%d bytes)\n\n", maxBytesIP, maxBytes)

			// printStats(stats)
		}
	}()

	buffer := make([]byte, 2048)

	// For every received packet:
	for {
		// - Extract metadata
		n, _, _ := conn.ReadFromUDP(buffer)

		var info PacketInfo
		// - Update statistics
		json.Unmarshal(buffer[:n], &info)

		update(stats, info.SrcIP, info.Length)
		update(stats, info.DstIP, info.Length)
	}
}
