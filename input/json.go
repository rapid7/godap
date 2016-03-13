package input

import (
  "bufio"
  "encoding/json"
  "github.com/rapid7/godap/api"
  "github.com/rapid7/godap/factory"
  "io"
)

type InputJson struct {
  scanner *bufio.Scanner
  FileSource
}

func (js *InputJson) ReadRecord() (data map[string]interface{}, err error) {
  if !js.scanner.Scan() {
    return nil, io.EOF
  }
  return data, json.Unmarshal(js.scanner.Bytes(), &data)
}

func NewInputJson(args []string) (input api.Input, err error) {
  inputJson := &InputJson{}
  var file string
  if len(args) > 0 {
    file = args[0]
  }
  err = inputJson.Open(file)
  inputJson.scanner = bufio.NewScanner(inputJson.fd)
  return inputJson, err
}

func init() {
  factory.RegisterInput("json", func(args []string) (input api.Input, err error) {
    return NewInputJson(args)
  })
}
