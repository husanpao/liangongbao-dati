package main

import (
	"bytes"
	"gitee.com/ekukuku/httpproxy"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func OnError(ctx *httpproxy.Context, where string,
	err *httpproxy.Error, opErr error) {
	// Log errors.
	log.Printf("ERR: %s: %s [%s]", where, err, opErr)
}

func OnAccept(ctx *httpproxy.Context, w http.ResponseWriter,
	req *http.Request) bool {
	// Handle local request has path "/info"
	//log.Printf("INFO:OnAccept Proxy: %s %s", req.Method, req.URL.String())
	if req.Method == "GET" && !req.URL.IsAbs() && req.URL.Path == "/info" {
		w.Write([]byte("This is go-httpproxy."))
		return true
	}
	return false
}

func OnAuth(ctx *httpproxy.Context, authType string, user string, pass string) bool {
	// Auth test user.
	if user == "admin" && pass == "jsepc@01!" {
		return true
	}
	return false
}

func OnConnect(ctx *httpproxy.Context, host string) (
	ConnectAction httpproxy.ConnectAction, newHost string) {
	// Apply "Man in the Middle" to all ssl connections. Never change host.
	return httpproxy.ConnectMitm, host
}

func OnRequest(ctx *httpproxy.Context, req *http.Request) (
	resp *http.Response) {
	// Log proxying requests.
	if strings.Index(req.URL.String(), "static") == -1 {
		log.Printf("INFO: Proxy: %s %s", req.Method, req.URL.String())
	}
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("OnRequest error %s", err.Error())
		return
	}
	if strings.Index(req.URL.String(), "ques/answerQues") != -1 {
		log.Printf("request Before:%s", string(s))
		requst := string(s)
		for {
			if strings.Index(requst, "(正确)") != -1 {
				requst = strings.Replace(requst, "(正确)", "", 1)
				req.ContentLength = req.ContentLength - 8
			} else {
				break
			}
		}
		s = []byte(requst)
		log.Printf("request After:%s", string(s))
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(s))

	return
}

func OnResponse(ctx *httpproxy.Context, req *http.Request,
	resp *http.Response) {
	// Add header "Via: go-httpproxy".
	if resp.Header.Get("Content-Type") != "application/json" || resp.Header.Get("Content-Encoding") != "gzip" {
		return
	}

	body := openGzip(resp.Body)
	log.Printf("resp Before:%s", body)
	r := Req{
		Url:  req.URL.String(),
		Data: body,
	}
	temp := updateBody(r)
	if temp != "" {
		body = temp
	}
	log.Printf("resp After:%s", body)
	bs := closeGzip(body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
}

func main() {

	// Create a new proxy with default certificate pair.
	prx, _ := httpproxy.NewProxy()

	// Set handlers.
	prx.OnError = OnError
	prx.OnAccept = OnAccept
	//prx.OnAuth = OnAuth
	prx.OnConnect = OnConnect
	prx.OnRequest = OnRequest
	prx.OnResponse = OnResponse
	log.Printf("%s", "Server start success on :35246 .")
	// Listen...
	http.ListenAndServe(":35246", prx)
}
