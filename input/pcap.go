// +build libpcap

package input

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"strconv"
	"strings"
	"time"
)

type InputPcap struct {
	handle *pcap.Handle
	opts   map[string]string
}

func (pcap *InputPcap) ReadRecord() (data map[string]interface{}, err error) {
	pktdata, ci, err := pcap.handle.ReadPacketData()
	packet := gopacket.NewPacket(pktdata, layers.LinkTypeEthernet, gopacket.Default)
	payload := packet.Data()
	packetSrc := ""
	packetDst := ""

	networkLayer := packet.NetworkLayer()
	if networkLayer != nil {
		networkFlow := networkLayer.NetworkFlow()
		packetSrc = networkFlow.Src().String()
		packetDst = networkFlow.Dst().String()
	}

	transportLayer := packet.TransportLayer()
	if transportLayer != nil {
		layerPayload := transportLayer.LayerPayload()
		if layerPayload != nil {
			payload = layerPayload
		}
	}

	return map[string]interface{}{
		"packet.src":       packetSrc,
		"packet.dst":       packetDst,
		"packet.data":      payload,
		"packet.timestamp": ci.Timestamp.UTC()}, err
}

func (pcap *InputPcap) ParseOpts(args []string) {
	for _, arg := range args {
		params := strings.SplitN(arg, "=", 2)
		if len(params) > 1 {
			pcap.opts[params[0]] = params[1]
		} else {
			pcap.opts[params[0]] = ""
		}
	}
}

func init() {
	factory.RegisterInput("pcap", func(args []string) (input api.Input, err error) {
		inputPcap := &InputPcap{}
		inputPcap.opts = map[string]string{
			"snaplen": "65536",
			"promisc": "false",
			"timeout": "-1",
			"rfmon":   "false",
		}
		inputPcap.ParseOpts(args)
		if file, ok := inputPcap.opts["file"]; ok {
			inputPcap.handle, err = pcap.OpenOffline(file)
		} else if iface, ok := inputPcap.opts["iface"]; ok {
			promisc, err := strconv.ParseBool(inputPcap.opts["promisc"])
			if err != nil {
				return nil, fmt.Errorf("Invalid promisc value: %s", inputPcap.opts["promisc"])
			}
			timeout, err := strconv.ParseInt(inputPcap.opts["timeout"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("Invalid timeout value: %s", inputPcap.opts["timeout"])
			}
			snaplen, err := strconv.ParseInt(inputPcap.opts["snaplen"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("Invalid snaplen value: %s", inputPcap.opts["snaplen"])
			}
			rfmon, err := strconv.ParseBool(inputPcap.opts["rfmon"])
			if err != nil {
				return nil, fmt.Errorf("Invalid rfmon value: %s", inputPcap.opts["rfmon"])
			}
			inactiveHandle, err := pcap.NewInactiveHandle(iface)
			if err != nil {
				return nil, fmt.Errorf("Could not create new inactive handle")
			}
			inactiveHandle.SetSnapLen(int(snaplen))
			inactiveHandle.SetPromisc(promisc)
			inactiveHandle.SetTimeout(time.Duration(timeout))
			inactiveHandle.SetRFMon(rfmon)
			inputPcap.handle, err = inactiveHandle.Activate()
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Either the file or iface option must be provided")
		}

		if v, ok := inputPcap.opts["filter"]; ok {
			bpferr := inputPcap.handle.SetBPFFilter(v)
			if bpferr != nil {
				err = fmt.Errorf("bpf filter complation failed: %s", bpferr)
			}
		}

		return inputPcap, err
	})
}
