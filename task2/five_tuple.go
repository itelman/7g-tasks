package task2

// FiveTuple is a representation of the 5-Tuple metadata of a packet
type FiveTuple struct {
	SrcIP    string
	DstIP    string
	SrcPort  uint16
	DstPort  uint16
	Protocol string
	Length   int
}
