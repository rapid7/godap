package api

type Input interface {
	ReadRecord() (record map[string]interface{}, err error)
}
