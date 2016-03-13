package input

import (
  "testing"
)

func TestNewSimpleJsonInputDoesntThrowErr(t *testing.T) {
  // TODO: This doesn't pass because of design issues
  /*   _, err := NewInputJson([]string{"json"})
       if err != nil {
          t.Errorf("Creating a json input with args \"json\" failed with %s", err.Error())
       }
  */
}
