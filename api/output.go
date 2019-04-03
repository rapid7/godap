package api

type Output interface {
	WriteRecord(doc map[string]interface{}) (err error)
	Start()
	Stop()
}
