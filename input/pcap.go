package input

import (
  "fmt"
  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
  "github.com/google/gopacket/pcap"
  "github.com/rapid7/godap/api"
  "github.com/rapid7/godap/factory"
  "strings"
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
  return map[string]interface{}{"packet.data": payload, "packet.timestamp": ci.Timestamp.UTC()}, err
}

func (pcap *InputPcap) ParseOpts(args []string) {
  pcap.opts = make(map[string]string)
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
    var file string
    if len(args) < 1 {
      panic("pcap input requires a filename argument")
    }
    file = args[0]
    if len(args) > 1 {
      inputPcap.ParseOpts(args)
    }
    inputPcap.handle, err = pcap.OpenOffline(file)
    if v, ok := inputPcap.opts["filter"]; ok {
      bpferr := inputPcap.handle.SetBPFFilter(v)
      if bpferr != nil {
        err = fmt.Errorf("bpf filter complation failed: %s", bpferr)
      }
    }
    return inputPcap, err
  })
}
