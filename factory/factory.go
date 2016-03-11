package factory

import (
   "errors"
   "fmt"
   "github.com/rapid7/godap/api"
   "sync"
)

var (
   inputsMu  sync.RWMutex
   outputsMu sync.RWMutex
   filtersMu sync.RWMutex
   inputs    = make(map[string]func(args []string) (api.Input, error))
   outputs   = make(map[string]func(args []string) (api.Output, error))
   filters   = make(map[string]func(args []string) (api.Filter, error))
)

func CreateInput(args []string) (input api.Input, err error) {
   plugin_name := args[0]
   if len(args) > 1 {
      args = args[1:]
   } else {
      args = []string{}
   }

   factory := inputs[plugin_name]
   if factory == nil {
      return nil, errors.New(fmt.Sprintf("Invalid input plugin: %s", plugin_name))
   }
   return factory(args)
}

func CreateOutput(args []string) (output api.Output, err error) {
   plugin_name := args[0]
   if len(args) > 1 {
      args = args[1:]
   } else {
      args = []string{}
   }

   factory := outputs[plugin_name]
   if factory == nil {
      return nil, errors.New(fmt.Sprintf("Invalid output plugin: %s", plugin_name))
   }
   return factory(args)
}

func CreateFilter(args []string) (filter api.Filter, err error) {
   plugin_name := args[0]
   if len(args) > 1 {
      args = args[1:]
   } else {
      args = []string{}
   }

   factory := filters[plugin_name]
   if factory == nil {
      return nil, errors.New(fmt.Sprintf("Invalid filter plugin: %s", plugin_name))
   }
   return factory(args)
}

func RegisterInput(name string, factory func(args []string) (api.Input, error)) {
   inputsMu.Lock()
   defer inputsMu.Unlock()
   inputs[name] = factory
}

func Inputs() []string {
   inputNames := []string{}
   for k := range inputs {
      inputNames = append(inputNames, k)
   }
   return inputNames
}

func RegisterOutput(name string, factory func(arg []string) (api.Output, error)) {
   outputsMu.Lock()
   defer outputsMu.Unlock()
   outputs[name] = factory
}

func Outputs() []string {
   outputNames := []string{}
   for k := range outputs {
      outputNames = append(outputNames, k)
   }
   return outputNames
}

func RegisterFilter(name string, factory func(arg []string) (api.Filter, error)) {
   filtersMu.Lock()
   defer filtersMu.Unlock()
   filters[name] = factory
}

func Filters() []string {
   filterNames := []string{}
   for k := range filters {
      filterNames = append(filterNames, k)
   }
   return filterNames
}
