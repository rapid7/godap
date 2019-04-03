package api

type Filter interface {
	Process(map[string]interface{}) (res []map[string]interface{}, err error)
}
