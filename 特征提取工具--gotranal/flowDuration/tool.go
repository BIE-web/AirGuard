package flowDuration

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func removeDuplicateElement(languages []int) []int {
	result := make([]int, 0, len(languages))
	temp := map[int]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return false
	}

	return ipv4[0] == 10 || // 10.0.0.0/8
		(ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31) || // 172.16.0.0/12
		(ipv4[0] == 169 && ipv4[1] == 254) || // 169.254.0.0/16
		(ipv4[0] == 192 && ipv4[1] == 168) // 192.168.0.0/16
}

// 源地址判定规则：
// 1. 一个外网地址，一个内网地址， 整内网地址为 源地址
// 2. 两个都是内网地址 地址大的为 源地址 例 10.1.1.1 于 192.168.123.5 则 192.168.123.5 为源地址
// 3. 两个都是外网地址 端口大的为源地址
// 4. ip相同时 端口大的为源地址
func dependSrc(src string, sport int, dst string, dport int) bool {
	ip1 := net.ParseIP(src)
	ip2 := net.ParseIP(dst)
	if HasLocalIP(ip1) {
		if HasLocalIP(ip2) {
			if ip1[0] < ip2[0] {
				return false
			} else if ip1[0] == ip2[0] {
				if ip1[1] < ip2[1] {
					return false
				} else if ip1[1] == ip2[1] {
					if ip1[2] < ip2[2] {
						return false
					} else if ip1[2] == ip2[2] {
						if ip1[3] < ip2[3] {
							return false
						} else if ip1[3] == ip2[3] {
							if ip1[4] < ip2[4] {
								return false
							} else if ip1[4] == ip2[4] {
								if sport < dport {
									return false
								}
							}
						}
					}
				}
			}
			return true

		} else {
			return true
		}
	} else if HasLocalIP(ip2) {
		return false
	} else {
		if sport >= dport {
			return true
		} else {
			return true
		}
	}
}

// 提取关键信息
// PCAP包长度 Timezone IPv4协议首部字节 TCPorUDP首部字节
// 目前只支撑IPV4

func getPacketInfo(packet *gopacket.Packet) (*[5]string, *Spacket) {
	timestamp := (*packet).Metadata().Timestamp
	packetSize := uint16((*packet).Metadata().Length)
	ipLayer := (*packet).Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		l3Byte := ipLayer.LayerContents()
		ip, _ := ipLayer.(*layers.IPv4)
		if ip.Protocol == layers.IPProtocolTCP {
			tcpLayer := (*packet).Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				l4Byte := tcpLayer.LayerContents()
				tcp, _ := tcpLayer.(*layers.TCP)
				dir := dependSrc(ip.SrcIP.String(), int(tcp.SrcPort), ip.DstIP.String(), int(tcp.DstPort))
				spacket := NewSpacket(dir, packetSize, timestamp, l3Byte, l4Byte)
				if dir {
					return &[5]string{ip.SrcIP.String(), strconv.Itoa(int(tcp.SrcPort)), ip.DstIP.String(), strconv.Itoa(int(tcp.DstPort)), "TCP"}, spacket
				} else {
					return &[5]string{ip.DstIP.String(), strconv.Itoa(int(tcp.DstPort)), ip.SrcIP.String(), strconv.Itoa(int(tcp.SrcPort)), "TCP"}, spacket
				}
			}
		} else if ip.Protocol == layers.IPProtocolUDP {
			udpLayer := (*packet).Layer(layers.LayerTypeUDP)
			if udpLayer != nil {
				l4Byte := udpLayer.LayerContents()
				udp, _ := udpLayer.(*layers.UDP)
				dir := dependSrc(ip.SrcIP.String(), int(udp.SrcPort), ip.DstIP.String(), int(udp.DstPort))
				spacket := NewSpacket(dir, packetSize, timestamp, l3Byte, l4Byte)
				if dir {
					return &[5]string{ip.SrcIP.String(), strconv.Itoa(int(udp.SrcPort)), ip.DstIP.String(), strconv.Itoa(int(udp.DstPort)), "UDP"}, spacket
				} else {
					return &[5]string{ip.DstIP.String(), strconv.Itoa(int(udp.DstPort)), ip.SrcIP.String(), strconv.Itoa(int(udp.SrcPort)), "UDP"}, spacket
				}
			}
		}
	}
	return nil, nil
}

func CalcMd5(packetInfo *[5]string) string {
	MD5 := md5.New()
	_, _ = io.WriteString(MD5, packetInfo[0]+" "+packetInfo[1]+" "+packetInfo[2]+" "+packetInfo[3]+" "+packetInfo[4])
	return hex.EncodeToString(MD5.Sum(nil))
}

func WriterCSV(data [][]string, path string) {

	File, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer File.Close()

	//创建写入接口
	writerCsv := csv.NewWriter(File)
	for _, line := range data {
		err = writerCsv.Write(line)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	writerCsv.Flush() //刷新，不刷新是无法写入的
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Send(server string, data [][]string) {
	lineData := ""
	for _, line := range data {
		for _, value := range line {
			lineData += value + " "
		}
		lineData = strings.TrimRight(lineData, " ")
		lineData += "\n"
	}
	lineData = strings.TrimRight(lineData, "\n")
	serverAddr, _ := net.ResolveUDPAddr("udp", server)
	conn, _ := net.DialUDP("udp", nil, serverAddr)
	_, err := conn.Write(String2Bytes(lineData))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	conn.Close()
}
