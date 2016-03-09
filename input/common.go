package input

import (
   "os"
   "github.com/rapid7/godap/util"
)

type FileSource struct {
   fd *os.File
}

func (fs *FileSource) Open(file_name string) error {
   if (util.StringInSlice(file_name, []string{ "", "-", "stdin" })) {
      fs.fd = os.Stdin   
   } else {
      fd, err := os.Open(file_name)
      if (err != nil) {
         return err
      }
      fs.fd = fd
   }
   return nil
}

func (fs *FileSource) Close() error {
   if (fs.fd != nil) {
      return fs.fd.Close()
   }
   return nil
}
