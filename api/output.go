package api

type Output interface {
   WriteRecord(doc map[string]interface{})
   Start()
   Stop()
}
