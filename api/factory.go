package api

type Factory interface {
   CreateInput(args string) (*Input, error)
   CreateOutput(args string) (*Output, error)
   CreateFilter(args string) (*Filter, error)
}
