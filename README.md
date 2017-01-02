# Simple Go HTTP client
[![Build Status](https://travis-ci.org/hooklift/httpclient.svg?branch=master)](https://travis-ci.org/hooklift/httpclient)
[![GoDoc](https://godoc.org/github.com/hooklift/httpclient?status.svg)](https://godoc.org/github.com/hooklift/httpclient)

Go's default HTTP client is well suited for most use cases but not for all, especially not for those that need a
long lived HTTP connection. For example:

* When downloading files using slow internet connections without abruptly interrupting it due to absolute deadlines kicking in
* When you want to open a long lived HTTP stream for receiving any sort of events or logs

Also, it is a common practice to globally share and tweak the default HTTP client and transport, causing issues that are
really difficult and time consuming to find. One specific example is adding or removing the client's deadline timeout.
If the former, valid and active connections will be interrupted. If the latter, and the server doesn't write anything back,
the connection blocks for as long as the operating system decides to time out, usually ~2-3 minutes.

This library will only compile with Go 1.7 or greater.

## Features

* Encourages to create a new HTTP client instance for each specific usage.
* Allows to set read/write resetable timeout on the underlined TCP connection.
* Remains context aware and connections can be canceled if the passed context is.


## Example

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hooklift/httpclient"
)

func main() {
	ExampleDialContext()
	ExampleDefault()
}

func ExampleDialContext() {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext:           httpclient.DialContext(5*time.Second, 5*time.Second),
			Proxy:                 http.ProxyFromEnvironment,
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
```
