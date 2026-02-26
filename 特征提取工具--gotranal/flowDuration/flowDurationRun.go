package flowDuration

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/gopacket"
)

func OnRun(packetSource *gopacket.PacketSource, duration int, outputname string, server string) {

	// flows[md5]
	flows := make(map[string]*Flow)
	packets := []*gopacket.Packet{}
	timeout := time.After(time.Duration(duration) * time.Second)

timeTicker:
	for {
		select {
		case <-timeout:
			break timeTicker
		case packet := <-packetSource.Packets():
			packets = append(packets, &packet)
		}
	}

	for _, packet := range packets {
		packetInfo, spacket := getPacketInfo(packet)
		if packetInfo != nil && spacket != nil {
			md5 := CalcMd5(packetInfo)
			packetTime := (*packet).Metadata().Timestamp.Truncate(time.Duration(duration) * time.Second)
			eachFlow, ok := flows[md5]
			if !ok {
				flows[md5] = NewFlow(packetInfo, packetTime, duration)
				flows[md5].addPacket(spacket)
			} else {
				eachFlow.addPacket(spacket)
			}
		}
	}
	if len(flows) == 0 {
		fmt.Printf("No traffic recv\n")
	} else {
		fmt.Printf("There are %d flows\n", len(flows))
		packetFeatureList := [][]string{}

		var wg sync.WaitGroup
		var lock sync.Mutex
		jobsChan := make(chan int, 20)
		packetFeatureList = append(packetFeatureList, getFeaturName())

		for _, eachFlow := range flows {
			wg.Add(1)
			jobsChan <- 1
			go func(eachFlow *Flow) {
				defer wg.Done()
				lock.Lock()
				defer lock.Unlock()
				packetFeatureList = append(packetFeatureList, eachFlow.getFeature(&flows))
				<-jobsChan
			}(eachFlow)

			wg.Wait()
		}
		Send(server, packetFeatureList[1:])
		//WriterCSV(packetFeatureList, outputname)
		fmt.Printf("Output path: %s\n", outputname)
	}
}

func OffRun(packetSource *gopacket.PacketSource, duration int, outputname string) {

	// flows[timestamp][md5]
	flows := make(map[int64]map[string]*Flow)
	//Read pcap
	for packet := range packetSource.Packets() {
		packetInfo, spacket := getPacketInfo(&packet)
		if packetInfo != nil && spacket != nil {
			//truncate函数取整时间段
			packetTime := packet.Metadata().Timestamp.Truncate(time.Duration(duration) * time.Second)
			eachDuration, ok := flows[packetTime.Unix()]
			md5 := CalcMd5(packetInfo)
			if !ok {
				eachDuration = make(map[string]*Flow)
				eachDuration[md5] = NewFlow(packetInfo, packetTime, duration)
				eachDuration[md5].addPacket(spacket)
				flows[packetTime.Unix()] = eachDuration
			} else {
				eachFlow, ok := eachDuration[md5]
				if !ok {
					eachDuration[md5] = NewFlow(packetInfo, packetTime, duration)
					eachDuration[md5].addPacket(spacket)
				} else {
					eachFlow.addPacket(spacket)
				}
			}
		}
	}
	packetFeatureList := [][]string{}

	var wg sync.WaitGroup
	var lock sync.Mutex
	jobsChan := make(chan int, 20)
	packetFeatureList = append(packetFeatureList, getFeaturName())
	cnt := 0

	for _, eachDuration := range flows {
		for _, eachFlow := range eachDuration {
			cnt += 1
			wg.Add(1)
			jobsChan <- 1
			go func(eachFlow *Flow) {
				defer wg.Done()
				lock.Lock()
				defer lock.Unlock()
				packetFeatureList = append(packetFeatureList, eachFlow.getFeature(&eachDuration))
				<-jobsChan
			}(eachFlow)

			wg.Wait()
		}
	}

	fmt.Printf("There are %d flows\n", cnt)
	WriterCSV(packetFeatureList, outputname)
	fmt.Printf("Output path: %s\n", outputname)
}
