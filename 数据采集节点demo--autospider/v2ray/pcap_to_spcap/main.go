package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	inputpath  string
	outputpath string
)

func init() {
	flag.StringVar(&inputpath, "i", "aaa", "usage for a")
	flag.StringVar(&outputpath, "o", ".", "usage for b")
}

func uint32_to_byte(v uint32) []byte {
	b := []byte{}
	b = append(b, uint8(v>>24))
	b = append(b, uint8(v>>16))
	b = append(b, uint8(v>>8))
	b = append(b, uint8(v))
	return b
}
func uint16_to_byte(v uint16) []byte {
	b := []byte{}
	b = append(b, uint8(v>>8))
	b = append(b, uint8(v))
	return b
}

func main() {
	start := time.Now()
	flag.Parse()

	handle, err := pcap.OpenOffline(inputpath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = handle.SetBPFFilter("tcp or udp")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer handle.Close()
	fp, err := os.Create(outputpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()
	buf := new(bytes.Buffer)
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			if ip.Protocol == layers.IPProtocolUDP {
				udpLayer := packet.Layer(layers.LayerTypeUDP)
				if udpLayer != nil {
					packetBin := uint16_to_byte(uint16(packet.Metadata().Length))
					packetBin = append(packetBin, uint32_to_byte(uint32(packet.Metadata().Timestamp.Unix()))...)
					packetBin = append(packetBin, uint32_to_byte(uint32(packet.Metadata().Timestamp.UnixMicro()-packet.Metadata().Timestamp.Unix()*1e6))...)
					packetBin = append(packetBin, ipLayer.LayerContents()...)
					packetBin = append(packetBin, udpLayer.LayerContents()...)
					packetBin = append([]byte{uint8(len(packetBin) + 1)}, packetBin...)
					binary.Write(buf, binary.LittleEndian, packetBin)
				}
			} else if ip.Protocol == layers.IPProtocolTCP {
				tcpLayer := packet.Layer(layers.LayerTypeTCP)
				if tcpLayer != nil {
					packetBin := uint16_to_byte(uint16(packet.Metadata().Length))
					packetBin = append(packetBin, uint32_to_byte(uint32(packet.Metadata().Timestamp.Unix()))...)
					packetBin = append(packetBin, uint32_to_byte(uint32(packet.Metadata().Timestamp.UnixMicro()-packet.Metadata().Timestamp.Unix()*1e6))...)
					packetBin = append(packetBin, ipLayer.LayerContents()...)
					packetBin = append(packetBin, tcpLayer.LayerContents()...)
					packetBin = append([]byte{uint8(len(packetBin) + 1)}, packetBin...)
					binary.Write(buf, binary.LittleEndian, packetBin)
				}
			}
		}
	}
	fp.Write(buf.Bytes())
	cost := time.Since(start)
	fmt.Printf("cost=[%s]", cost)
}
