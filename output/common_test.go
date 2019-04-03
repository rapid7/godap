package output

import "testing"

const (
	INVALID_UTF8 = "\xc3\x28"
	VALID_UTF8   = "\xc3\xb1"
)

func TestSanitizeNonUTF8Characters(t *testing.T) {
	fileDestination := new(FileDestination)
	if v := fileDestination.Sanitize(INVALID_UTF8); v != "(" {
		t.Errorf("Expected \"(\", received %s", v)
	}
}

func TestSanitizeValidUTF8Characters(t *testing.T) {
	fileDestination := new(FileDestination)
	if v := fileDestination.Sanitize(VALID_UTF8); v != "ñ" {
		t.Errorf("Expected ñ, received %s", v)
	}
}

func TestSanitizeMap(t *testing.T) {
	fileDestination := new(FileDestination)
	hash := fileDestination.Sanitize(map[string]interface{}{INVALID_UTF8: INVALID_UTF8})
	if _, ok := hash.(map[string]interface{})["("]; !ok {
		t.Errorf("Expected \"(\" map key to exist, map is %s", hash)
	}
}

func TestSanitizeSlice(t *testing.T) {
	fileDestination := new(FileDestination)
	data := []string{INVALID_UTF8, VALID_UTF8}
	arr := fileDestination.Sanitize(data).([]string)
	if arr[0] != "(" || arr[1] != VALID_UTF8 {
		t.Errorf("Expected \"(\" array with one element, array is %s and of size %d", arr, len(arr))
	}
}
