package go_request

import (
	"fmt"
	"testing"
	"time"
)

func TestGo(t *testing.T) {

	res, err := Go(&Request{
		Url:    "http://api.tansoz.cn/15.php",
		Method: "GET",
		Async:  true,
		Data: map[string]interface{}{
			"id": "65535",
		},
		DataType: "json",
		//ContentType: "multipart/form-data",
		//ContentType: "application/json",
		Error: func(err error) {
			fmt.Println("async err:", err)
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
