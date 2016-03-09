package output

import (
   "fmt"
   "github.com/rapid7/godap/api"
   "github.com/rapid7/godap/factory"
)

type OutputLines struct {
   // TODO: Need to port FileDestination here
}

func (lines *OutputLines) WriteRecord(data map[string]interface{}) {
   // TODO
   fmt.Println(data)
}

func (lines *OutputLines) Start() {
}

func (lines *OutputLines) Stop() {
}

func init() {
   factory.RegisterOutput("lines", func(args []string) (api.Output, error) {
      var lines *OutputLines = &OutputLines{}
      // TODO
      return lines, nil
   });
}
