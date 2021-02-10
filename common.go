package go_request

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"reflect"
)

const (
	VERSION = "1.0.0" // build version. usage update version. bug fix version.
)

func MD5(s string) string {
	b := md5.Sum([]byte(s))
	return hex.EncodeToString(b[0:])
}

func Base64_Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func json_encode(data interface{}) string {

	if b, err := json.Marshal(data); err == nil {
		return string(b)
	} else {
		panic(err)
	}
	return ""
}

func json_decode(jdata []byte, i interface{}) interface{} {

	var p interface{}

	if i != nil {
		p = reflect.New(reflect.TypeOf(i)).Interface()
	}

	if err := json.Unmarshal(jdata, &p); err != nil {
		panic(err)
	}
	return p
}
