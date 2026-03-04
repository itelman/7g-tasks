package task2

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

/*
Objective:

Observes traffic involving a given machine/IP (example: 192.168.10.100), packets that come from/to
the machine, and delivers basic info to the analyzer

Steps:

1. Capture and read packets, parse layers.
2. Extracts: 5-tuple, packet length
3. Sends JSON metadata via UDP to analyzer.
*/

func RunSniffer() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: ./sniffer <interface> <filter_ip> <analyzer_host:port>")
	}

	interf := os.Args[1]
	filterIP := os.Args[2]
	analyzerAddr := os.Args[3]

	// Capture all packets seen by the given interface
	handle, err := pcap.OpenLive(interf, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	err = handle.SetBPFFilter("host " + filterIP)
	if err != nil {
		log.Fatal(err)
	}

	// Creates a continuous stream of captured packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	conn, err := net.Dial("udp", analyzerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Sniffer in progress...")

	// Graceful shutdown for termination with Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Printf("\nSniffer terminated.\n")
		os.Exit(0)
	}()

	for packet := range packetSource.Packets() {
		// Process IPv4 packets.
		// If no IP layer -> skip.
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}

		ip := ipLayer.(*layers.IPv4)

		// Extract 5-Tuple (SrcIP, DstIP, SrcPort, DstPort, Protocol)
		info := PacketInfo{
			SrcIP:  ip.SrcIP.String(),
			DstIP:  ip.DstIP.String(),
			Length: len(packet.Data()),
		}

		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)
			info.SrcPort = uint16(tcp.SrcPort)
			info.DstPort = uint16(tcp.DstPort)
			info.Proto = "TCP"
		} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp := udpLayer.(*layers.UDP)
			info.SrcPort = uint16(udp.SrcPort)
			info.DstPort = uint16(udp.DstPort)
			info.Proto = "UDP"
		}

		// Send metadata to analyzer via UDP
		data, _ := json.Marshal(info)

		// For debugging
		// fmt.Println(info)

		conn.Write(data)
	}
}
