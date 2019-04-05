package input

import (
	"errors"
	"github.com/dgraph-io/badger"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"github.com/rapid7/godap/util"
	"strconv"
)

type InputBadger struct {
	db              *badger.DB
	txn             *badger.Txn
	iterator        *badger.Iterator
	prefetch_values bool
	keys_only       bool
	prefix          []byte
}

func (ib *InputBadger) ReadRecord() (data map[string]interface{}, err error) {
	if ib.prefix != nil {
		if !ib.iterator.ValidForPrefix(ib.prefix) {
			return nil, errors.New("EOF")
		}
	} else if !ib.iterator.Valid() {
		return nil, errors.New("EOF")
	}

	data = make(map[string]interface{})
	item := ib.iterator.Item()
	if val, err := item.ValueCopy(nil); err == nil {
		data["key"] = string(item.Key())
		if !ib.keys_only {
			data["value"] = string(val)
		}
	}
	ib.iterator.Next()
	return data, err // TODO error
}

func (ib *InputBadger) Start() {
	ib.txn = ib.db.NewTransaction(false)
	iterator_opts := badger.DefaultIteratorOptions
	iterator_opts.PrefetchValues = ib.prefetch_values // TODO: allow customization
	ib.iterator = ib.txn.NewIterator(iterator_opts)
	if ib.prefix != nil {
		ib.iterator.Seek(ib.prefix)
	} else {
		ib.iterator.Rewind()
	}
}

func (ib *InputBadger) Stop() {
	ib.iterator.Close()
	ib.txn.Discard() // TODO: Error handling
	ib.db.Close()
}

func init() {
	factory.RegisterInput("badger", func(args []string) (lines api.Input, err error) {
		ib := &InputBadger{}
		parsed_opts := util.ParseOpts(args)
		opts := badger.DefaultOptions
		if opts.Dir = parsed_opts["dir"]; opts.Dir == "" {
			return nil, errors.New("Invalid or no badger database `dir` directory specified")
		}
		if opts.ValueDir = parsed_opts["value_dir"]; opts.ValueDir == "" {
			opts.ValueDir = opts.Dir
		}
		if val_str, ok := parsed_opts["prefetch_values"]; ok {
			if ib.prefetch_values, err = strconv.ParseBool(val_str); err != nil {
				return nil, err
			}
		}
		if val_str, ok := parsed_opts["keys_only"]; ok {
			if ib.keys_only, err = strconv.ParseBool(val_str); err != nil {
				return nil, err
			}
		}
		if parsed_opts["prefix"] != "" {
			ib.prefix = []byte(parsed_opts["prefix"])
		}
		if ib.db, err = badger.Open(opts); err != nil {
			return nil, err // TODO more graceful handling
		}
		return ib, err
	})
}
