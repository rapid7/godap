package util

func Merge(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		dest[k] = v
	}
	return dest
}
