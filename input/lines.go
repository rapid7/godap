package input

import (
   "bufio"
   "github.com/rapid7/godap/api"
   "github.com/rapid7/godap/factory"
)

type InputLines struct {
   reader *bufio.Reader
   FileSource
}

func (lines *InputLines) ReadRecord() (data map[string]interface{}, err error) {
   line, err := lines.reader.ReadString('\n')
   if (line != "") {
      line = line[:len(line)-1]
   }
   return map[string]interface{} { "line": line }, err
}

func init() {
   factory.RegisterInput("lines", func(args []string) (input api.Input, err error) {
      inputLines := &InputLines{}
      var file string
      if (len(args) > 0) {
         file = args[0]
      }
      err = inputLines.Open(file)
      inputLines.reader = bufio.NewReader(inputLines.fd)
      return inputLines, err
   });
}
