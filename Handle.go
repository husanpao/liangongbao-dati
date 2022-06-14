package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Req struct {
	Url  string `json:"url"`
	Data string `json:"data"`
}

func openGzip(body io.ReadCloser) string {
	data, err := gzip.NewReader(body)
	if err != nil {
		log.Printf("error gzip.NewReader %s", err.Error())
		return ""
	}
	s, err := ioutil.ReadAll(data)
	if err != nil {
		log.Printf("error ioutil.ReadAll %s", err.Error())
		return ""
	}
	return string(s)
}
func closeGzip(body string) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(body)); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return b.Bytes()
}
func updateBody(req Req) string {
	res, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	request, _ := http.NewRequest("POST", "http://81.68.160.189:35248/handle", bytes.NewReader(res))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("post data error:%v\n", err)
	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return string(respBody)
	}
	return ""
}
