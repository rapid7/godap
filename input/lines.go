package input

import (
   "bufio"
   "errors"
   "github.com/rapid7/godap/api"
   "github.com/rapid7/godap/factory"
)

type InputLines struct {
   scanner *bufio.Scanner
   FileSource
}

func (lines *InputLines) ReadRecord() (data map[string]interface{}, err error) {
   if !lines.scanner.Scan() {
      return nil, errors.New("eof")
   }
   return map[string]interface{}{"line": lines.scanner.Text()}, err
}

func init() {
   factory.RegisterInput("lines", func(args []string) (input api.Input, err error) {
      inputLines := &InputLines{}
      var file string
      if len(args) > 0 {
         file = args[0]
      }
      err = inputLines.Open(file)
      inputLines.scanner = bufio.NewScanner(inputLines.fd)
      return inputLines, err
   })
}
