package go_request

import (
	"errors"
	"net/http"
)

type Response struct {
	HttpResponse *http.Response
	Data         interface{}
}

func New(request *Request) (res *Response, err error) {

	if request.Async {

		go func() {

			defer func() {

				// catch errors
				if tmp := recover(); tmp != nil {

					switch t := tmp.(type) {

					case string:
						err = errors.New(t)
					default:
						err = tmp.(error)
					}

					if request.Error != nil {
						request.Error(err)
					}
				}
			}()

			request.open()
			request.send()

			resp := request.read()

			if request.Complete != nil {
				request.Complete(&Response{
					HttpResponse: resp,
					Data:         nil,
				})
			}

		}()
		return nil, nil
	} else {

		defer func() {

			// catch errors
			if tmp := recover(); tmp != nil {

				switch t := tmp.(type) {

				case string:
					err = errors.New(t)
				default:
					err = tmp.(error)
				}
			}
		}()

		request.open()
		request.send()

		resp := request.read()

		return &Response{
			HttpResponse: resp,
			Data:         nil,
		}, err
	}
}
