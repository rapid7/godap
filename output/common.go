package output

import (
	"github.com/rapid7/godap/util"
	"os"
	"unicode/utf8"
)

type FileDestination struct {
	fd *os.File
}

func (filedest *FileDestination) Open(file_name string) error {
	filedest.Close()
	if util.StringInSlice(file_name, []string{"", "-", "stdout"}) {
		filedest.fd = os.Stdout
	} else {
		fd, err := os.Create(file_name)
		if err != nil {
			return err
		}
		filedest.fd = fd
	}
	return nil
}

func (filedest *FileDestination) Close() {
	if filedest.fd != nil {
		filedest.fd.Close()
	}
	filedest.fd = nil
}

func (filedest *FileDestination) Sanitize(o interface{}) interface{} {
	if s, ok := o.(string); ok {
		if !utf8.ValidString(s) {
			v := make([]rune, 0, len(s))
			for i, r := range s {
				if r == utf8.RuneError {
					_, size := utf8.DecodeRuneInString(s[i:])
					if size == 1 {
						continue
					}
				}
				v = append(v, r)
			}
			s = string(v)
		}
		return s
	} else if v, ok := o.(map[string]interface{}); ok {
		r := make(map[string]interface{})
		for key, val := range v {
			if safekey, ok := filedest.Sanitize(key).(string); ok {
				r[safekey] = filedest.Sanitize(val)
			}
		}
		return r
	} else if v, ok := o.([]string); ok {
		r := make([]string, 0, len(v))
		for _, item := range v {
			r = append(r, filedest.Sanitize(item).(string))
		}
		return r
	}

	return o
}
