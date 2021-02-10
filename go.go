package go_request

import (
	"errors"
	"net/http"
)

type Response struct {
	HttpResponse *http.Response
	Data         interface{}
}

func Go(request *Request) (res *Response, err error) {

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

			if resp.Data != nil && request.Success != nil {
				request.Success(resp)
			}

			if request.Complete != nil {
				request.Complete(resp)
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

		return request.read(), err

	}
}
