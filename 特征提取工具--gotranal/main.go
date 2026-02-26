package main

import (
	"flag"
	"fmt"
	"gotranal/flowDuration"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	ifacename string
	bpf       string
	server    string

	inputpath     string
	isTimeoutMode bool
	timeout       time.Duration

	outputname string
	duration   time.Duration

	err    error
	handle *pcap.Handle
)

type OnLineOfflineFlagSet struct {
	*flag.FlagSet
	cmdComment string
}

func main() {
	onlineCmd := &OnLineOfflineFlagSet{
		FlagSet:    flag.NewFlagSet("on", flag.ExitOnError),
		cmdComment: "Online mode",
	}
	onlineCmd.StringVar(&outputname, "o", "", "outputpath prefix (default {YYYYMMDDHHMMSS}.csv)")
	onlineCmd.StringVar(&ifacename, "I", "any", "iface name")
	onlineCmd.StringVar(&bpf, "b", "tcp or udp", "bpf filter")
	onlineCmd.StringVar(&server, "s", "127.0.0.1:31115", "vpn_finder server addr:port")
	onlineCmd.DurationVar(&duration, "d", time.Second*60, "sniff timeout")

	offlineCmd := &OnLineOfflineFlagSet{
		FlagSet:    flag.NewFlagSet("off", flag.ExitOnError),
		cmdComment: "Offline mode",
	}
	offlineCmd.StringVar(&inputpath, "i", "", "Input path")
	offlineCmd.StringVar(&outputname, "o", "", "Output name")
	offlineCmd.BoolVar(&isTimeoutMode, "m", false, "Drution_mode or timeout_mode. Recommend drution_mode")
	offlineCmd.DurationVar(&timeout, "t", time.Second*600, "Use in timeout_mode")
	offlineCmd.DurationVar(&duration, "d", time.Second*60, "Use in duration_mode")

	subcommands := map[string]*OnLineOfflineFlagSet{
		onlineCmd.Name():  onlineCmd,
		offlineCmd.Name(): offlineCmd,
	}

	useage := func() {
		fmt.Printf("Usage: gotranal COMMAND\n\n")
		for _, v := range subcommands {
			fmt.Printf("%s %s\n", v.Name(), v.cmdComment)
			v.PrintDefaults()
			fmt.Println()
		}
		os.Exit(2)
	}

	if len(os.Args) < 2 { // 即没有输入子命令
		useage()
	}

	cmd := subcommands[os.Args[1]] // 第二个参数必须是我们支持的子命令
	if cmd == nil {
		useage()
	}

	cmd.Parse(os.Args[2:])

	if os.Args[1] == "on" {
		for {
			name := outputname + time.Now().Format("20060102150405") + ".csv"
			handle, err = pcap.OpenLive(ifacename, 65535, true, pcap.BlockForever)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			err = handle.SetBPFFilter(bpf)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			defer handle.Close()
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			flowDuration.OnRun(packetSource, int(duration)/1e9, name, server)
		}
	} else {
		start := time.Now()
		if inputpath == "" {
			fmt.Println("require inputpath")
			os.Exit(1)
		}
		if outputname == "" {
			outputname = inputpath + ".csv"
		}

		handle, err = pcap.OpenOffline(inputpath)

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		// 只对tcp和udp流进行提取
		err = handle.SetBPFFilter("tcp or udp")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer handle.Close()

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		if isTimeoutMode {

		} else {
			flowDuration.OffRun(packetSource, int(duration)/1e9, outputname)
		}

		cost := time.Since(start)
		fmt.Printf("Time cost: %s\n", cost)
	}
}
