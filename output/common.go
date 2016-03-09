package output

import (
   "fmt"
   "os"
   "github.com/rapid7/godap/util"
)

type FileDestination struct {
   fd *os.File
}

func (filedest *FileDestination) Open(file_name string) error {
   filedest.Close()
   if (util.StringInSlice(file_name, []string{ "", "-", "stdin" })) {
      fmt.Println("Opening stdout...")
      filedest.fd = os.Stdout   
   } else {
      fmt.Printf("Opening %s...\n", file_name)
      fd, err := os.Open(file_name)
      if (err != nil) {
         return err
      }
      filedest.fd = fd
   }
   return nil
}

func (filedest *FileDestination) Close() {
   if (filedest.fd != nil) {
      filedest.fd.Close()
   }
   filedest.fd = nil
}

func (filedest *FileDestination) Sanitize(o interface{}) interface{} {
   if v, ok := o.(string); ok {
      // TODO: Encode?
      return v
   } else if v, ok := o.(map[string]interface{}); ok {
      r := make(map[string]interface{})
      for key, val := range v {
         if safekey, ok := filedest.Sanitize(key).(string); ok {
            r[safekey] = filedest.Sanitize(val)
         }
      }
      return r
   } else if v, ok := o.([]interface{}); ok {
      r := make([]interface{}, len(v), len(v))
      return r
   }

   return nil
}
