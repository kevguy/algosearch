package web

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/dimfeld/httptreemux/v5"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}

// Query returns the web call queries from the request.
// https://blog.csdn.net/quicmous/article/details/81322015
func Query(r *http.Request, key string) []string {
	u, _ := url.Parse(r.URL.String())
	values, _ := url.ParseQuery(u.RawQuery)
	//fmt.Println(u)           // /time?a=111&b=1212424
	//fmt.Println(u.RawQuery)  // a=111&b=1212424
	//fmt.Println(values)      // map[a:[111] b:[1212424]]
	//fmt.Println(values["a"]) //[111]
	//fmt.Println(values["b"]) //[1212424]
	return values[key]
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
