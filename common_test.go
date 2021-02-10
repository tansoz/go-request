package go_request

import (
	"fmt"
	"testing"
)

func Test_json_decode(t *testing.T) {

	a := json_decode([]byte(`{"name":"asdads","age":15}`), struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{})
	fmt.Println(a)
}
