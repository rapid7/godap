package filter

import (
   "fmt"
   "github.com/rapid7/godap/api"
   "github.com/rapid7/godap/factory"
   "strings"
)

/////////////////////////////////////////////////
// select filter
/////////////////////////////////////////////////
type FilterSelect struct {
   Base
}

func (fs *FilterSelect) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   ndoc := make(map[string]interface{})
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok {
         ndoc[k] = docv
      }
   }
   ndocs := make([]map[string]interface{}, 0)
   if len(ndoc) > 0 {
      ndocs = append(ndocs, ndoc)
   }
   return ndocs, nil
}

func init() {
   factory.RegisterFilter("select", func(args []string) (lines api.Filter, err error) {
      filterSelect := &FilterSelect{}
      filterSelect.ParseOpts(args)
      return filterSelect, nil
   })
}

/////////////////////////////////////////////////
// rename filter
/////////////////////////////////////////////////

type FilterRename struct {
   Base
}

func (fs *FilterRename) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      if _, ok := doc[k]; ok {
         doc[v] = doc[k]
         delete(doc, k)
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("rename", func(args []string) (lines api.Filter, err error) {
      filterRename := &FilterRename{}
      filterRename.ParseOpts(args)
      return filterRename, nil
   })
}

/////////////////////////////////////////////////
// remove filter
/////////////////////////////////////////////////

type FilterRemove struct {
   Base
}

func (fs *FilterRemove) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if _, ok := doc[k]; ok {
         delete(doc, k)
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("remove", func(args []string) (lines api.Filter, err error) {
      filterRemove := &FilterRemove{}
      filterRemove.ParseOpts(args)
      return filterRemove, nil
   })
}

/////////////////////////////////////////////////
// insert filter
/////////////////////////////////////////////////

type FilterInsert struct {
   Base
}

func (fs *FilterInsert) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      doc[k] = v
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("insert", func(args []string) (lines api.Filter, err error) {
      filterInsert := &FilterInsert{}
      filterInsert.ParseOpts(args)
      return filterInsert, nil
   })
}

/////////////////////////////////////////////////
// include filter
/////////////////////////////////////////////////

type FilterInclude struct {
   Base
}

func (fs *FilterInclude) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      if docv, ok := doc[k]; ok && strings.Contains(docv.(string), v) {
         return []map[string]interface{}{doc}, nil
      }
   }
   return make([]map[string]interface{}, 0), nil
}

func init() {
   factory.RegisterFilter("include", func(args []string) (lines api.Filter, err error) {
      filterInclude := &FilterInclude{}
      filterInclude.ParseOpts(args)
      return filterInclude, nil
   })
}

/////////////////////////////////////////////////
// exclude filter
/////////////////////////////////////////////////

type FilterExclude struct {
   Base
}

func (fs *FilterExclude) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, v := range fs.opts {
      if docv, ok := doc[k]; ok && strings.Contains(docv.(string), v) {
         return make([]map[string]interface{}, 0), nil
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("exclude", func(args []string) (lines api.Filter, err error) {
      filterExclude := &FilterExclude{}
      filterExclude.ParseOpts(args)
      return filterExclude, nil
   })
}

/////////////////////////////////////////////////
// exists filter
/////////////////////////////////////////////////

type FilterExists struct {
   Base
}

func (fs *FilterExists) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok && len(docv.(string)) > 0 {
         return []map[string]interface{}{doc}, nil
      }
   }
   return make([]map[string]interface{}, 0), nil
}

func init() {
   factory.RegisterFilter("exists", func(args []string) (lines api.Filter, err error) {
      filterExists := &FilterExists{}
      filterExists.ParseOpts(args)
      return filterExists, nil
   })
}

/////////////////////////////////////////////////
// notexists filter
/////////////////////////////////////////////////

type FilterNotExists struct {
   Base
}

func (fs *FilterNotExists) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   for k, _ := range fs.opts {
      if docv, ok := doc[k]; ok && len(docv.(string)) > 0 {
         return make([]map[string]interface{}, 0), nil
      }
   }
   return []map[string]interface{}{doc}, nil
}

func init() {
   factory.RegisterFilter("notexists", func(args []string) (lines api.Filter, err error) {
      filterNotExists := &FilterNotExists{}
      filterNotExists.ParseOpts(args)
      return filterNotExists, nil
   })
}

/////////////////////////////////////////////////
// where filter
/////////////////////////////////////////////////

type FilterWhere struct {
   key      string
   operator func(string, string) bool
   value    string
   Base
}

func (fs *FilterWhere) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
   if docv, ok := doc[fs.key]; ok && fs.operator(fs.value, docv.(string)) {
      return []map[string]interface{}{doc}, nil
   }
   return make([]map[string]interface{}, 0), nil
}

func init() {
   factory.RegisterFilter("where", func(args []string) (lines api.Filter, err error) {
      filterWhere := &FilterWhere{}
      if len(args) != 3 {
         panic(fmt.Sprintf("Expected 3 arguments to 'where' but got %d: %s", len(args), args))
      }
      filterWhere.key = args[0]
      filterWhere.value = args[2]
      if args[1] == "==" {
         filterWhere.operator = func(lhs string, rhs string) bool {
            return lhs == rhs
         }
      } else if args[1] == "!=" {
         filterWhere.operator = func(lhs string, rhs string) bool {
            return lhs != rhs
         }
      } else {
         panic(fmt.Sprintf("Unknown conditional operator for 'where': %s", args[1]))
      }
      return filterWhere, nil
   })
}
