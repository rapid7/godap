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
  transportLayer := packet.TransportLayer()
  if transportLayer != nil {
    layerPayload := transportLayer.LayerPayload()
    if layerPayload != nil {
      payload = layerPayload
    }
  }
  networkLayer := packet.NetworkLayer()
  packetSrc := ""
  packetDst := ""
  if networkLayer != nil {
    networkFlow := networkLayer.NetworkFlow()
    packetSrc = networkFlow.Src().String()
    packetDst = networkFlow.Dst().String()
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
      inputPcap.handle, err = pcap.OpenLive(iface, int32(snaplen), promisc, time.Duration(timeout))
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
