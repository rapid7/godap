package output

import (
   "strings"
   "regexp"
   "github.com/rapid7/godap/api"
   "github.com/rapid7/godap/factory"
   "github.com/rapid7/godap/util"
)

const FIELD_WILDCARD = "_"

type OutputLines struct {
   delimiter string
   fields []string
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

   lines.fd.WriteString(strings.Join(out, lines.delimiter))
   lines.fd.WriteString("\n")
   lines.fd.Sync()
}

func (lines *OutputLines) Start() {
}

func (lines *OutputLines) Stop() {
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
      if (header && !util.StringInSlice(outputLines.delimiter, outputLines.fields)) {
         outputLines.fd.WriteString(strings.Join(outputLines.fields, outputLines.delimiter) + "\n")
         outputLines.fd.Sync()
      }
      return outputLines, nil
   });
}
