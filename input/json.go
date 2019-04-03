package input

import (
	"bufio"
	"encoding/json"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
)

type InputJson struct {
	reader *bufio.Reader
	FileSource
}

func (js *InputJson) ReadRecord() (data map[string]interface{}, err error) {
	text, err := js.reader.ReadString('\n')
	if text != "" {
		text = text[:len(text)-1]
	}
	return data, json.Unmarshal([]byte(text), &data)
}

func NewInputJson(args []string) (input api.Input, err error) {
	inputJson := &InputJson{}
	var file string
	if len(args) > 0 {
		file = args[0]
	}
	err = inputJson.Open(file)
	inputJson.reader = bufio.NewReader(inputJson.fd)
	return inputJson, err
}

func init() {
	factory.RegisterInput("json", func(args []string) (input api.Input, err error) {
		return NewInputJson(args)
	})
}
