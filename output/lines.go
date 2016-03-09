package output

import (
   "strings"
   "regexp"
   "bufio"
   "github.com/rapid7/godap/api"
   "github.com/rapid7/godap/factory"
   "github.com/rapid7/godap/util"
)

const FIELD_WILDCARD = "_"

type OutputLines struct {
   delimiter string
   fields []string
   writer *bufio.Writer
   FileDestination
}

func (lines *OutputLines) WriteRecord(data map[string]interface{}) {
   var out []string

   if (util.StringInSlice(FIELD_WILDCARD, lines.fields)) {
      for _, v := range data {
         out = append(out, lines.Sanitize(v).(string))
      }
   } else {
      for _, field := range lines.fields {
         sanitized := lines.Sanitize(data[field])
         if (sanitized != nil) {
            out = append(out, sanitized.(string))
         }
      }
   }

   if (len(out) < 1) { return }

   lines.writer.WriteString(strings.Join(out, lines.delimiter))
   lines.writer.WriteString("\n")
   lines.writer.Flush()
}

func (lines *OutputLines) Start() {
}

func (lines *OutputLines) Stop() {
   lines.writer.Flush()
   lines.Close()
}

func init() {
   factory.RegisterOutput("lines", func(args []string) (lines api.Output, err error) {
      outputLines := &OutputLines{}
      var file string
      outputLines.delimiter = ","
      outputLines.fields = []string{ FIELD_WILDCARD }

      header := false

      re := regexp.MustCompile("(?i)^[ty1]")
      for _, arg := range args {
         params := strings.SplitN(arg, "=", 2)
         switch params[0] {
         case "file":
            file = params[1]
            break
         case "header":
            header = re.MatchString(params[1])
            break
         case "fields":
            outputLines.fields = strings.Split(params[1], ",")
            break
         case "delimiter":
            switch params[1] {
            case "tab":
               outputLines.delimiter = "\t"
               break
            case "null":
               outputLines.delimiter = "\x00"
               break
            default:
               outputLines.delimiter = params[1]
            }
         }
      }

      err = outputLines.Open(file)
      outputLines.writer = bufio.NewWriter(outputLines.fd)
      if (header && !util.StringInSlice(FIELD_WILDCARD, outputLines.fields)) {
         outputLines.writer.WriteString(strings.Join(outputLines.fields, outputLines.delimiter) + "\n")
         outputLines.writer.Flush()
      }
      return outputLines, nil
   });
}
