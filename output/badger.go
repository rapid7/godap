package output

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"github.com/rapid7/godap/util"
	"strings"
)

type OutputBadger struct {
	db             *badger.DB
	batch          *badger.WriteBatch
	string_builder *strings.Builder
	key_field      string
	value_field    string
}

func (ob *OutputBadger) WriteRecord(data map[string]interface{}) (err error) {
	key := []byte(data[ob.key_field].(string))
	value := []byte("")
	if val := data[ob.value_field]; val != nil {
		value = []byte(val.(string))
	}
	return ob.batch.Set(key, value, 0)
}

func (ob *OutputBadger) Start() {
	ob.batch = ob.db.NewWriteBatch()
}

func (ob *OutputBadger) Stop() {
	ob.batch.Flush() // TODO: Handle error
	ob.batch.Cancel()

	// Now run gc...
	fmt.Println("running gc..")
	for err := badger.ErrNoRewrite; err != badger.ErrNoRewrite; err = ob.db.RunValueLogGC(0.5) {
		fmt.Println("running gc again..")
	}

	defer ob.db.Close()
}

func init() {
	factory.RegisterOutput("badger", func(args []string) (lines api.Output, err error) {
		ob := &OutputBadger{}
		parsed_opts := util.ParseOpts(args)
		ob.string_builder = &strings.Builder{}
		opts := badger.DefaultOptions
		if ob.key_field = parsed_opts["key_field"]; ob.key_field == "" {
			return nil, errors.New("No `key_field` specified.")
		}
		ob.value_field = parsed_opts["value_field"]
		if opts.Dir = parsed_opts["dir"]; opts.Dir == "" {
			return nil, errors.New("Invalid or no badger database `dir` directory specified")
		}
		if opts.ValueDir = parsed_opts["value_dir"]; opts.ValueDir == "" {
			opts.ValueDir = opts.Dir
		}
		if ob.db, err = badger.Open(opts); err != nil {
			return nil, err
		}
		return ob, nil
	})
}
