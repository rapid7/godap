package filter

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"io/ioutil"
	"net/http"
)

/////////////////////////////////////////////////
// transform filter
/////////////////////////////////////////////////

type FilterIp struct {
	BaseFilter
}

func (fs *FilterIp) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
	for k := range fs.opts {
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(doc[k].([]byte))))
		if err != nil {
			doc[fmt.Sprintf("%s.ip", k)], _ = ioutil.ReadAll(req.Body)
		}
	}
	return []map[string]interface{}{doc}, nil
}

func init() {
	factory.RegisterFilter("ip", func(args []string) (lines api.Filter, err error) {
		FilterIp := &FilterIp{}
		FilterIp.ParseOpts(args)
		return FilterIp, nil
	})
}
