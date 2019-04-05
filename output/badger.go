package output

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"strings"
)

type OutputBadger struct {
	db             *badger.DB
	batch          *badger.WriteBatch
	string_builder *strings.Builder
}

func (ob *OutputBadger) WriteRecord(data map[string]interface{}) (err error) {
	// construct key
	// we could use a different method here, since keys only have to be unique
	// per value, and there is typically a very finite number of those.
	// However, for now we'll just be lazy and concatenate the value to make a
	// unique primary key...
	ob.string_builder.Reset()
	ob.string_builder.WriteString(data["name"].(string))
	ob.string_builder.WriteString(",")
	ob.string_builder.WriteString(data["type"].(string))
	ob.string_builder.WriteString(",")
	ob.string_builder.WriteString(data["value"].(string))
	var key = ob.string_builder.String()
	var value = data["value"].(string)
	return ob.batch.Set([]byte(key), []byte(value), 0)
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
		ob.string_builder = &strings.Builder{}
		opts := badger.DefaultOptions
		opts.Dir = "./badger-test"
		opts.ValueDir = "./badger-test"
		if ob.db, err = badger.Open(opts); err != nil {
			return nil, err // TODO more graceful handling
		}
		return ob, nil
	})
}
