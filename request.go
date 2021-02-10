package go_request

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type RWC interface { // mean Reader Writer Close
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)

	Close() error
}

type Request struct {
	Method      string                      // request method GET,POST,etc. default 'GET'
	Url         string                      // the request link.
	Async       bool                        // the request is use async way or not. default 'false'
	Host        string                      // the server host name.
	Timeout     int                         // if request time over the millisecond value with no response then callback the function. default '15000'
	Headers     map[string]string           // set request headers.
	DataType    string                      // auto convert receive data.	('json','xml')
	ContentType string                      // the request data content type. default 'application/x-www-form-urlencoded'
	Data        map[string]interface{}      // the request data.
	File        map[string]string           // the request file.
	Cookie      map[string]string           // cookies.
	Success     func(res *Response)         // if request is success then callback the function.
	Fail        func(res *Response)         // if request is fail then callback the function.
	Complete    func(res *Response)         // if request is complete then callback the function.
	Error       func(err error)             // when the request occur error then invoke the function.
	StatusCode  map[int]func(res *Response) // according to response status code to callback specific function.

	link          *url.URL // parsed url
	connection    RWC      // connection
	err           error    // keep error information
	contentLength int64    // the request body length
}

func (this *Request) defaultParams() {
	this.checkMethod()
	this.checkHeaders()
	this.checkTimeout()
	this.checkContentType()
	this.checkDataType()

	this.checkLink()
}

func (this *Request) checkMethod() {

	this.Method = strings.ToUpper(strings.TrimSpace(this.Method)) // upper method string and delete right or left blank char.

	switch this.Method {

	case "GET", "POST", "HEAD", "DELETE", "PUT", "OPTION":
		// do nothing
	default:
		this.Method = "GET" // set default value
	}
}
func (this *Request) checkHeaders() {
	// default data
	tmp := map[string]string{
		"user-agent": "Go-Request/" + VERSION + "; (+https://github.com/tansoz/go-request)",
	}

	// check
	for k, i := range this.Headers {

		tmp[strings.ToLower(k)] = i
	}

	this.Headers = tmp
}
func (this *Request) checkTimeout() {

	if this.Timeout == 0 {
		this.Timeout = 15000
	}
}

func (this *Request) checkContentType() {

	t := strings.TrimSpace(strings.ToLower(this.ContentType)) // remove left and right space and make the string lower

	switch t {

	// supported content type
	case "application/x-www-form-urlencoded",
		"multipart/form-data",
		"application/json":
		this.ContentType = t

	default:
		this.ContentType = "application/x-www-form-urlencoded" // set default value to content type

	}
}

func (this *Request) checkLink() {

	var err error

	if this.Url == "" {
		panic("The URL can't be empty.")
	}

	this.link, err = url.Parse(this.Url)
	if err != nil {
		panic(err)
	}
}

func (this *Request) checkDataType() {

	var err error

	this.DataType = strings.TrimSpace(strings.ToLower(this.DataType))

	switch this.DataType {

	case "json": // supported data type
		// do noting
	default:
		this.DataType = ""

	}

	this.link, err = url.Parse(this.Url)
	if err != nil {
		panic(err)
	}
}

func (this *Request) open() {

	this.defaultParams() // check default value

	// init connection, trying to connect the server.
	port := this.link.Port()
	switch this.link.Scheme {

	case "http":
		if port == "" {
			port = ":80" // default http port
		} else {
			port = ""
		}

		if this.connection, this.err = net.Dial("tcp", this.link.Host+port); this.err != nil {
			panic(this.err)
		}
	case "https":
		if port == "" {
			port = ":443" // default https port
		} else {
			port = ""
		}

		if this.connection, this.err = tls.Dial("tcp", this.link.Host+port, &tls.Config{
			InsecureSkipVerify: true, // not verify certificate
		}); this.err != nil {
			panic(this.err)
		}
	case "test": // debugger
		if this.connection, this.err = os.OpenFile(MD5(this.link.Host), os.O_CREATE|os.O_RDWR, 777); this.err != nil {
			panic(this.err)
		}
	}
}

func (this *Request) getBoundary() string {

	return "----WebKitFormBoundary" + Base64_Encode(this.Url + "Go-Request/" + VERSION)[0:16]
}

func escape(s string) string {
	return regexp.MustCompile("([\"\\\\])").ReplaceAllString(s, "\\$1")
}

func (this *Request) Write(data []byte) {

	if _, err := this.connection.Write(data); err != nil {
		panic(err)
	}
}

func http_query(data interface{}) string {

	s := ""

	switch t := data.(type) {

	case map[string]interface{}:
		for k, v := range t {
			s += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(fmt.Sprint(v))
		}
	case map[string]string:
		for k, v := range t {
			s += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(v)
		}
	default:
		s = "&" // do noting
	}
	return s[1:]
}

func (this *Request) send() {

	if this.connection != nil {

		fps := make(map[string]*os.File)
		query := ""                    // query data
		requestHeader := ""            // http request header
		boundary := this.getBoundary() // POST form data boundary
		formDataBody := ""             // POST the data of form data

		if this.File != nil {

			for k, i := range this.File {
				if fps[k], this.err = os.Open(i); this.err != nil {
					panic(this.err)
				}
			}
		}

		if "multipart/form-data" == this.ContentType {

			if this.Method == "POST" {
				this.Headers["content-type"] = "multipart/form-data; boundary=" + boundary

				if this.Data != nil {

					for k, i := range this.Data {
						formDataBody += boundary + "\r\nContent-Disposition: form-data; name=\"" + escape(k) + "\"\r\n\r\n" + fmt.Sprint(i) + "\r\n"
					}

					this.contentLength = int64(len(formDataBody)) + 40 // 40 is last boundary and ending '--'
				}

				for k, fp := range fps {

					this.contentLength += 138
					this.contentLength += int64(len(escape(k)))
					if stat, err := fp.Stat(); err == nil {
						this.contentLength += stat.Size()
						this.contentLength += int64(len(escape(fp.Name())))
					} else {
						panic(err)
					}
				}
			}

		} else if "application/json" == this.ContentType {

			if this.Method == "POST" {

				for k, fp := range fps {
					this.Data[k] = fp.Name()
				}
				query = json_encode(this.Data)
				this.contentLength = int64(len(query))
			}

		} else {
			// file to query
			for k, i := range fps {
				query += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(i.Name())
			}

			if this.Data != nil {

				for k, i := range this.Data {
					query += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(fmt.Sprint(i))
				}
			}

			if len(query) > 0 {

				query = query[1:] // remove first '&'

				if this.Method != "POST" {
					if this.link.RawQuery != "" {
						this.link.RawQuery = this.link.Query().Encode() + "&" + query
					} else {
						this.link.RawQuery = query
					}
				} else {
					this.Headers["content-type"] = this.ContentType
					this.contentLength = int64(len(query))
				}
			}

		}

		this.Headers["content-length"] = fmt.Sprint(this.contentLength)

		if this.Host == "" {
			this.Headers["host"] = this.link.Host
		} else {
			this.Headers["host"] = this.Host
		}

		headers := ""
		for k, i := range this.Headers {

			headers += k + ": " + i + "\r\n"
		}

		requestHeader = fmt.Sprintf("%s %s HTTP/1.0\r\n%s\r\n", this.Method, this.link.RequestURI(), headers)

		// send request header
		this.Write([]byte(requestHeader))

		// send request body
		if this.Method == "POST" {
			if query != "" {
				// send query data
				this.Write([]byte(query))
			} else if formDataBody != "" {
				// send form data
				this.Write([]byte(formDataBody))

				// send file if have
				for k, fp := range fps {

					// send file info
					this.Write([]byte(boundary + "\r\nContent-Disposition: form-data; name=\"" + escape(k) + "\"; filename=\"" + escape(fp.Name()) + "\"\r\nContent-Type: application/octet-stream\r\n\r\n"))
					// send file data
					b := make([]byte, 1024)
					for {

						if rn, err := fp.Read(b); err != nil && err != io.EOF {
							panic(err)
						} else {
							this.Write(b[0:rn])
							if rn <= 0 {
								this.Write([]byte("\r\n"))
								break
							}
						}
					}
				}

				// send the end of the body
				this.Write([]byte(boundary + "--"))
			}
		}

	} else {
		panic("Call the open function before calling the send function.")
	}
}

func (this *Request) read() *Response {

	if resp, err := http.ReadResponse(bufio.NewReader(this.connection), nil); err == nil {

		switch this.DataType {

		case "json":
			var tmp bytes.Buffer
			if _, err := tmp.ReadFrom(resp.Body); err != nil {
				panic(err)
			} else {
				return &Response{
					HttpResponse: resp,
					Data:         json_decode(tmp.Bytes(), nil),
				}
			}

		}

		return &Response{
			HttpResponse: resp,
			Data:         nil,
		}

	} else {
		panic(err)
	}
	return nil
}
