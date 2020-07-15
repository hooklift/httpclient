package httpclient_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hooklift/httpclient"
)

func ExampleDialContext() {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext:           httpclient.DialContext(5*time.Second, 5*time.Second),
			Proxy:                 http.ProxyFromEnvironment,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
	res, err := client.Get("https://google.com")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body[:]))
}

func ExampleDefault() {
	client := httpclient.Default() // read/write timeout: 30s, connect timeout: 10s
	res, err := client.Get("https://google.com")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body[:]))
}
