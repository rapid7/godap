package api

type Filter interface {
   run() bool
}
