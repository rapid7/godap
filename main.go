package main

import (
  "github.com/mattn/go-shellwords"
  "github.com/rapid7/godap/api"
  "github.com/rapid7/godap/factory"
  _ "github.com/rapid7/godap/filter"
  _ "github.com/rapid7/godap/input"
  _ "github.com/rapid7/godap/output"
  "log"
  "os"
  "regexp"
  "runtime"
  "strings"
  "sync"
)

const VERSION = "0.0.1"

func main() {
  Console := log.New(os.Stderr, "", 0)
  var args [][]string
  trace := false

  re := regexp.MustCompile("\\s*\\+\\s*")
  if len(os.Args) > 1 {
    for _, bit := range re.Split(strings.Join(os.Args[1:], " "), -1) {
      aset, _ := shellwords.Parse(bit)
      if len(aset) < 1 {
        usage(Console)
      }

      arg := aset[0]

      if arg == "--trace" {
        trace = true
        arg, aset = aset[0], aset[1:]
      }

      if arg == "-h" || arg == "--help" {
        usage(Console)
      }

      if arg == "--version" || arg == "-v" {
        version(Console)
      }

      if arg == "--inputs" {
        show_inputs(Console)
      }

      if arg == "--outputs" {
        show_outputs(Console)
      }

      if arg == "--filters" {
        show_filters(Console)
      }

      args = append(args, aset)
    }
  }

  if len(args) < 2 {
    usage(Console)
  }

  inp_args, args := args[0], args[1:]
  out_args, args := args[len(args)-1], args[:len(args)-1]

  inp, err := factory.CreateInput(inp_args)
  if err != nil {
    Console.Printf("Error: %s", err)
    usage(Console)
  }
  out, err := factory.CreateOutput(out_args)
  if err != nil {
    Console.Printf("Error: %s", err)
    usage(Console)
  }

  filters := []api.Filter{}
  for _, arg := range args {
    filter, err := factory.CreateFilter(arg)
    if err != nil {
      Console.Printf("Error: %s", err)
      usage(Console)
    }
    filters = append(filters, filter)
  }

  out.Start()

  numcpu := runtime.NumCPU()
  runtime.GOMAXPROCS(numcpu)
  inch := make(chan map[string]interface{}, 1000)
  go func() {
    for {
      doc, err := inp.ReadRecord()
      if err != nil {
        break
      }
      if doc == nil {
        continue
      }
      inch <- doc
    }
    close(inch)
  }()

  var wg sync.WaitGroup
  outch := make(chan map[string]interface{})
  wg.Add(numcpu)
  for i := 0; i < numcpu; i++ {
    go func() {
      defer wg.Done()
      for doc := range inch {
        docs := []map[string]interface{}{doc}
        for _, filter := range filters {
          for _, doc := range docs {
            docs, err = filter.Process(doc)
          }
        }
        for _, doc := range docs {
          outch <- doc
        }
      }
    }()
  }

  var wg2 sync.WaitGroup
  wg2.Add(1)
  go func() {
    defer wg2.Done()
    for data := range outch {
      err := out.WriteRecord(data)
      if err != nil {
        panic(err)
      }
    }
  }()
  wg.Wait()
  close(outch)
  wg2.Wait()

  if trace {
    Console.Println("shouldn't see this")
  }
  out.Stop()
}

func usage(Console *log.Logger) {
  Console.Println("")
  Console.Printf("  Usage: %s [input] + [filter] + [output]\n", os.Args[0])
  Console.Println("       --inputs")
  Console.Println("       --outputs")
  Console.Println("       --filters")
  Console.Println("")
  Console.Printf("Example: echo world | %s lines stdin + rename line=hello + json stdout\n", os.Args[0])
  Console.Println("")
  os.Exit(1)
}

func version(Console *log.Logger) {
  Console.Printf("dap %s", VERSION)
  os.Exit(1)
}

func show_inputs(console *log.Logger) {
  console.Println("Inputs:")
  for _, k := range factory.Inputs() {
    console.Printf(" * %s", k)
  }
  os.Exit(1)
}

func show_outputs(console *log.Logger) {
  console.Println("Outputs:")
  for _, k := range factory.Outputs() {
    console.Printf(" * %s", k)
  }
  os.Exit(1)
}

func show_filters(console *log.Logger) {
  console.Println("Filters:")
  for _, k := range factory.Filters() {
    console.Printf(" * %s", k)
  }
  os.Exit(1)
}
