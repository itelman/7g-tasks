package task2

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Stats struct {
	Packets int
	Bytes   int
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
			fmt.Printf("Top by bytes: %s (%d bytes)\n", maxBytesIP, maxBytes)
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

// Each packet updates stats
// Time complexity: O(1)
func update(stats map[string]*Stats, ip string, length int) {
	if _, ok := stats[ip]; !ok {
		stats[ip] = &Stats{}
	}
	stats[ip].Packets++
	stats[ip].Bytes += length
}
