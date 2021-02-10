## 灵感
> 我遇到想上传文件的需求，但是 Golang 自带的 HTTP 库我弄了很久也不懂到底要怎么弄。而且平常也觉得自带的库，好是好，可是我觉得还是不够简单。也因为我平常写前端代码，用到 AJAX，所以我觉得我自己应该自己封装一个用法类似 AJAX，用起来更方便的 HTTP 请求库。然后我就写下了这一个库。

# Go-Request

> A HTTP request library by Golang.  
> The usage like JQuery AJAX.

Usage:

### Install
```
go get -u github.com/tansoz/go-request
```

### Demo
```go
package main

import . "github.com/tansoz/go-request"

func main(){
	Go(&Request{
		Method:"GET",
		Url:"https://www.github.com",
		Async:true,
		Success:func(res *Response){
			// when request success callback this function. 
		},
		Complete:func(res *Response){
			// when request have not error will callback this function. 
		},
		Error:func(err error){
            
		},
	})
}
```

# Go-请求库

> 一个 HTTP 的 Golang（Go语言）网络请求库。  
> 用法上面像 JQuery 的 AJAX 。

用法:

### 安装
```
go get -u github.com/tansoz/go-request
```