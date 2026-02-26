package flowDuration

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
)

func calStatistics(arr []float64) (float64, float64, float64, float64, float64, float64, float64, float64, float64, []float64) {
	var mean, min, max, iqr, q1, q2, q3, std, variance float64
	modes := []float64{}
	mean, min, max, iqr, q1, q2, q3, std, variance = 0, 0, 0, 0, 0, 0, 0, 0, 0
	if len(arr) < 1 {
		return mean, min, max, iqr, q1, q2, q3, std, variance, modes
	} else {
		mean, _ = stats.Mean(arr)
		mean, _ = stats.Round(mean, 6)

		min, _ = stats.Min(arr)
		min, _ = stats.Round(min, 6)

		max, _ = stats.Max(arr)
		max, _ = stats.Round(max, 6)

		quartiles, _ := stats.Quartile(arr)
		q1, _ = stats.Round(quartiles.Q1, 6)
		q2, _ = stats.Round(quartiles.Q2, 6)
		q3, _ = stats.Round(quartiles.Q3, 6)

		iqr, _ = stats.InterQuartileRange(arr)
		iqr, _ = stats.Round(iqr, 6)

		std, _ = stats.StandardDeviation(arr)
		std, _ = stats.Round(std, 6)

		variance, _ = stats.Variance(arr)
		variance, _ = stats.Round(variance, 6)

		modes, _ = stats.Mode(arr)
		return mean, min, max, iqr, q1, q2, q3, std, variance, modes
	}
}

func getFeaturName() []string {
	return []string{
		"srcIP", "srcPort", "dstIP", "dstPort", "l4Proto",
		"connSipDip", "connSipDprt",
		"aveIat", "minIat", "maxIat", "iqrIat", "q1Iat", "q2Iat", "q3Iat", "stdIat", "varIat",
		"avePktSz", "minPktSz", "maxPktSz", "iqrPktSz", "q1PktSz", "q2PktSz", "q3PktSz", "stdPktSz", "varPktSz", "modePktSz",
		"bytps", "pktps",
		"bytAsm", "pktAsm",
		"ipToS",
		"ipMinTTL", "ipMaxTTL", "ipTTLChg ",
		"ipAbnormalLenth",
		"l4AbnormalLenth",
		"tcpPSeqCnt", "tcpPAckCnt",
		"tcpAveWinSz", "tcpMinWinSz", "tcpMaxWinSz", "tcpWinSzDwnCnt", "tcpWinSzUpCnt", "tcpWinSzChgDirCnt",
		"tcpFlag",
	}
}

type Flow struct {
	src      string
	sport    string
	dst      string
	dport    string
	protocol string
	//startTime 用于六元组模式下，以drution时间段为分流判断依据
	startTime  time.Time
	duration   int
	packets    []*Spacket
	fwdPackets []*Spacket
	bwdPackets []*Spacket
}

func NewFlow(packeInfo *[5]string, timestamp time.Time, duration int) *Flow {
	return &Flow{
		src:        packeInfo[0],
		sport:      packeInfo[1],
		dst:        packeInfo[2],
		dport:      packeInfo[3],
		protocol:   packeInfo[4],
		startTime:  timestamp,
		duration:   duration,
		packets:    []*Spacket{},
		fwdPackets: []*Spacket{},
		bwdPackets: []*Spacket{},
	}
}

func (f *Flow) String() string {
	return fmt.Sprintf("%s:%s <-> %s:%s %s packets:%d time:%d", f.src, f.sport, f.dst, f.dport, f.protocol, f.getPacketsNum(), f.startTime.Unix())
}

func (f *Flow) addPacket(packet *Spacket) {
	f.packets = append(f.packets, packet)
}

func (f *Flow) dividePackets() {
	for _, eachPacket := range f.packets {
		if eachPacket.dir {
			f.fwdPackets = append(f.fwdPackets, eachPacket)
		} else {
			f.bwdPackets = append(f.bwdPackets, eachPacket)
		}
	}
}

func (f *Flow) getPacketsNum() int {
	return len(f.packets)
}

// 流间特征 统计同一内部IP、同一时间段连接的目的IP或相同端口的个数
func (f *Flow) getConn(flows *map[string]*Flow) []float64 {

	connSipNum := float64(0)
	connSipDipNum := float64(0)
	connSipDprtNum := float64(0)

	for _, v := range *flows {

		if v.src == f.src {
			if v.dst == f.dst {
				connSipDipNum += 1
			}
			if v.dport == f.dport {
				connSipDprtNum += 1
			}
			connSipNum += 1
		}
	}

	connSipDip := connSipDipNum / connSipNum
	connSipDprt := connSipDprtNum / connSipNum

	return []float64{connSipDip, connSipDprt}
}

// L3层特征 包大小和间隔到达时间统计 (平均值 最小值 最大值 四分位距 三个四等分点 标准差 方差 众数)
func (f *Flow) getPktIatSize() ([]float64, string) {
	var preTime, nextTime time.Time
	preTime = f.packets[0].timestamp
	packetIat := []float64{}
	packetSize := []float64{}
	for _, eachPacket := range f.packets {
		packetSize = append(packetSize, float64(eachPacket.packetSize))
		nextTime = eachPacket.timestamp
		packetIat = append(packetIat, math.Abs((float64(nextTime.UnixNano())-float64(preTime.UnixNano()))/1e9))
		preTime = nextTime
	}
	// 连续值不需要 众数
	packetIatMean, packetIatMin, packetIatMax, packetIatIqr, packetIatQ1, packetIatQ2, packetIatQ3, packetIatStd, packetIatVar, _ := calStatistics(packetIat[1:])
	packetSizeMean, packetSizeMin, packetSizeMax, packetSizeIqr, packetSizeQ1, packetSizeQ2, packetSizeQ3, packetSizeStd, packetSizeVar, packetSizeModeList := calStatistics(packetSize)
	packetSizeMode := strconv.Itoa(int(f.packets[0].packetSize))
	if len(packetSizeModeList) != 0 {
		packetSizeMode = strconv.Itoa(int(packetSizeModeList[0]))
	}
	return []float64{packetIatMean, packetIatMin, packetIatMax, packetIatIqr, packetIatQ1, packetIatQ2, packetIatQ3, packetIatStd, packetIatVar,
		packetSizeMean, packetSizeMin, packetSizeMax, packetSizeIqr, packetSizeQ1, packetSizeQ2, packetSizeQ3, packetSizeStd, packetSizeVar}, packetSizeMode
}

// L3层特征  每秒发字节数 每秒发包数
func (f *Flow) getPacketOrByteTps() []float64 {
	//六元组模式下流间隔为duration
	packetPs := float64(len(f.packets)) / float64(f.duration)
	packetSize := []float64{}

	for _, eachPacket := range f.packets {
		packetSize = append(packetSize, float64(eachPacket.packetSize))
	}
	packetSizeSum, _ := stats.Sum(packetSize)
	packetSizeSum, _ = stats.Round(packetSizeSum, 6)
	bytePs := packetSizeSum / (float64(f.duration))

	return []float64{bytePs, packetPs}
}

func (f *Flow) getPacketOrByteASM() []float64 {

	fwdFlowByte := 0
	allFlowByte := 0
	for _, eachPacket := range f.fwdPackets {
		fwdFlowByte += int(eachPacket.packetSize)
	}
	for _, eachPacket := range f.packets {
		allFlowByte += int(eachPacket.packetSize)
	}

	return []float64{float64(fwdFlowByte) / float64(allFlowByte), float64(len(f.fwdPackets)) / float64(len(f.packets))}
}

// L3层特征 根据T2 对每个包的Tos字段做或运算
func (f *Flow) getIPToS() int {
	tos := int(0)
	for _, eachPacket := range f.packets {
		tos |= eachPacket.IPToS()
	}
	return tos
}

// L3层特征 ipTTLchg TTL变化数/包数
func (f *Flow) getIPTTL() []float64 {
	firstIp := f.packets[0]
	minTTL := firstIp.IPTTL()
	maxTTL := firstIp.IPTTL()
	perTTL := firstIp.IPTTL()
	ipTTLChg := int(0)
	if len(f.packets) == 1 {
		return []float64{float64(minTTL), float64(maxTTL), float64(0)}
	}
	for _, eachPacket := range f.packets[1:] {
		nowTTL := eachPacket.IPTTL()
		if nowTTL > maxTTL {
			maxTTL = nowTTL
		}
		if nowTTL < minTTL {
			minTTL = nowTTL
		}
		if nowTTL != perTTL {
			ipTTLChg += 1
			perTTL = nowTTL
		}
	}
	return []float64{float64(minTTL), float64(maxTTL), float64(ipTTLChg) / float64(len(f.packets))}
}

// L3层特征 ipv4长度一般为20，长度不为20和整条流的比值
func (f *Flow) getIPAbnormalLenth() float64 {
	num := 0
	for _, eachPacket := range f.packets {
		if len(eachPacket.l3Byte) != 20 {
			num += 1
		}
	}
	return float64(num) / float64(f.getPacketsNum())
}

// L4层特征 TCP长度一般位20 ，长度不为20和整条流的比值， UDP不做记录
func (f *Flow) getL4AbnormalLenth() float64 {
	num := 0
	if f.protocol == "UDP" {
		return float64(-1)
	}
	for _, eachPacket := range f.packets {
		if len(eachPacket.l4Byte) != 20 {
			num += 1
		}
	}
	return float64(num) / float64(f.getPacketsNum())
}

// L4层特征 记录TCP的序列号和确认号变化，UDP不做记录
func (f *Flow) getTCPSeqAckCnt() []float64 {
	if f.protocol == "UDP" {
		return []float64{-1, -1}
	} else {
		seqList := []int{}
		ackList := []int{}
		for _, eachPacket := range f.packets {
			seqList = append(seqList, eachPacket.TCPSeq())
			ackList = append(ackList, eachPacket.TCPAck())
		}
		seqList = removeDuplicateElement(seqList)
		ackList = removeDuplicateElement(ackList)

		return []float64{float64(len(seqList)) / float64(f.getPacketsNum()), float64(len(ackList)) / float64(f.getPacketsNum())}
	}
}

// L4层特征 TCP窗口
func (f *Flow) getTCPSWinSz() []float64 {
	if f.protocol == "UDP" {
		return []float64{-1, -1, -1, -1, -1, -1}
	} else {
		perWinSz, minWinSz, maxWinSz, allWinSz, downWinCnt, upWinCnt, chgWinCnt := f.packets[0].TCPWin(), f.packets[0].TCPWin(), f.packets[0].TCPWin(), 0, 0, 0, 0
		for _, eachPacket := range f.packets {
			nowTCPWin := eachPacket.TCPWin()
			if perWinSz < nowTCPWin {
				chgWinCnt += 1
				upWinCnt += 1
				perWinSz = nowTCPWin
			} else if perWinSz > nowTCPWin {
				chgWinCnt += 1
				downWinCnt += 1
				perWinSz = nowTCPWin
			}
			if minWinSz > nowTCPWin {
				minWinSz = nowTCPWin
			}
			if maxWinSz < nowTCPWin {
				maxWinSz = nowTCPWin
			}
			allWinSz += nowTCPWin
		}

		return []float64{float64(allWinSz) / float64(f.getPacketsNum()), float64(minWinSz), float64(maxWinSz),
			float64(downWinCnt) / float64(f.getPacketsNum()), float64(upWinCnt) / float64(f.getPacketsNum()),
			float64(chgWinCnt) / float64(f.getPacketsNum())}
	}
}

// L4 TCPFlag
func (f *Flow) getTcpFlag() int {
	if f.protocol == "UDP" {
		return -1
	} else {
		flag := int(0)
		for _, eachPacket := range f.packets {
			flag |= eachPacket.TCPFlag()
		}
		return flag
	}
}

func (f *Flow) getFeature(flows *map[string]*Flow) []string {
	featureList := []string{f.src, f.sport, f.dst, f.dport, f.protocol}
	f.dividePackets()

	for _, v := range f.getConn(flows) {
		featureList = append(featureList, strconv.FormatFloat(v, 'f', 6, 64))
	}

	PacketIatSizeStats, PacketSizeMode := f.getPktIatSize()
	for _, v := range PacketIatSizeStats {
		featureList = append(featureList, strconv.FormatFloat(v, 'f', 6, 64))
	}
	featureList = append(featureList, PacketSizeMode)
	for _, v := range f.getPacketOrByteTps() {
		featureList = append(featureList, strconv.FormatFloat(v, 'f', 6, 64))
	}
	for _, v := range f.getPacketOrByteASM() {
		featureList = append(featureList, strconv.FormatFloat(v, 'f', 6, 64))
	}
	featureList = append(featureList, strconv.Itoa(f.getIPToS()))
	for _, v := range f.getIPTTL() {
		featureList = append(featureList, strconv.FormatFloat(v, 'f', 6, 64))
	}
	featureList = append(featureList, strconv.FormatFloat(f.getIPAbnormalLenth(), 'f', 6, 64))
	featureList = append(featureList, strconv.FormatFloat(f.getL4AbnormalLenth(), 'f', 6, 64))
	for _, v := range f.getTCPSeqAckCnt() {
		featureList = append(featureList, strconv.FormatFloat(v, 'f', 6, 64))
	}
	for _, v := range f.getTCPSWinSz() {
		featureList = append(featureList, strconv.FormatFloat(v, 'f', 6, 64))
	}
	featureList = append(featureList, strconv.Itoa(f.getTcpFlag()))
	return featureList
}
