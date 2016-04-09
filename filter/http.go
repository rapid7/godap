package filter

import ( 
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
	"bufio"
	"github.com/rapid7/godap/factory"
	"github.com/rapid7/godap/api"
)
/////////////////////////////////////////////////
// transform filter
/////////////////////////////////////////////////

type FilterHttp struct {
  BaseFilter
}

func (fs *FilterHttp) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
  for k, _ := range fs.opts {
	  req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader([]byte(doc[k].(string)))))
	  if err != nil && req != nil {
	  	if (req.Body != nil) {
	  	doc[fmt.Sprintf("%s.http.request.content", k)], _ = ioutil.ReadAll(req.Body)
	  }
	  }
	}
  return []map[string]interface{}{doc}, nil
}

func init() {
  factory.RegisterFilter("http", func(args []string) (lines api.Filter, err error) {
    filterHttp := &FilterHttp{}
    filterHttp.ParseOpts(args)
    return filterHttp, nil
  })
}