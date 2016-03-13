package output

import (
  "bufio"
  "encoding/json"
  "github.com/rapid7/godap/api"
  "github.com/rapid7/godap/factory"
)

type OutputJson struct {
  writer *bufio.Writer
  FileDestination
}

func (oj *OutputJson) WriteRecord(data map[string]interface{}) (err error) {
  json, err := json.Marshal(oj.Sanitize(data))
  if err == nil {
    _, err = oj.writer.Write(json)
    if err == nil {
      err = oj.writer.WriteByte('\n')
      if err == nil {
        err = oj.writer.Flush()
      }
    }
  }

  return err
}

func (oj *OutputJson) Start() {
}

func (oj *OutputJson) Stop() {
}

func init() {
  factory.RegisterOutput("json", func(args []string) (lines api.Output, err error) {
    var file string
    if len(args) > 0 {
      file = args[0]
    }
    outputJson := &OutputJson{}
    err = outputJson.Open(file)
    outputJson.writer = bufio.NewWriter(outputJson.fd)
    return outputJson, nil
  })
}
