package go_request

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	res, err := New(&Request{
		Url:    "https://www.baidu.com/",
		Method: "GET",
		Async:  true,
		//ContentType: "multipart/form-data",
		Error: func(err error) {
			fmt.Println("async", err)
		},
		Complete: func(res *Response) {
			fmt.Println("async res:", res)
		},
	})

	fmt.Println("sync res:", res)
	fmt.Println("sync err:", err)

	for {
		time.Sleep(1000000)
	}

}
