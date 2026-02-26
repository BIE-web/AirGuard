package flowDuration

import (
	"encoding/binary"
	"time"
)

// 分析加密流量无需Payload信息,为此专门设计Spacket,消去payload信息，保留IP协议首部， TCPorUDP首部， 新增加timezone，packetSize, dir字段
// dir 为流方向， 若于dependSrc 方向规定相同 则为true
type Spacket struct {
	dir        bool
	packetSize uint16
	timestamp  time.Time
	l3Byte     []byte
	l4Byte     []byte
}

func NewSpacket(dir bool, packetSize uint16, timestamp time.Time, l3Byte []byte, l4Byte []byte) *Spacket {
	if len(l3Byte) >= 20 && len(l4Byte) >= 8 {
		return &Spacket{
			dir:        dir,
			packetSize: packetSize,
			timestamp:  timestamp,
			l3Byte:     l3Byte,
			l4Byte:     l4Byte,
		}
	} else {
		return nil
	}
}

func (p *Spacket) IPToS() int {
	return int(p.l3Byte[1])
}

func (p *Spacket) IPTTL() int {
	return int(p.l3Byte[8])
}

func (p *Spacket) TCPSeq() int {
	if len(p.l3Byte) < 20 {
		return -1
	} else {
		Seq := p.l4Byte[4:8]
		for i := 0; i < len(Seq)/2; i++ {
			j := len(Seq) - i - 1
			Seq[i], Seq[j] = Seq[j], Seq[i]
		}
		return int(binary.LittleEndian.Uint32(Seq))
	}
}

func (p *Spacket) TCPAck() int {
	if len(p.l3Byte) < 20 {
		return -1
	} else {
		Ack := p.l4Byte[8:12]
		return int(binary.BigEndian.Uint32(Ack))
	}
}

func (p *Spacket) TCPWin() int {
	if len(p.l3Byte) < 20 {
		return -1
	} else {
		Win := p.l4Byte[14:16]
		return int(binary.BigEndian.Uint16(Win))
	}
}

func (p *Spacket) TCPFlag() int {
	if len(p.l3Byte) < 20 {
		return -1
	} else {
		return int(p.l4Byte[13])
	}
}
