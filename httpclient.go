// Package httpclient adds support for resetable read/write timeouts to Go's HTTP transport and exposes
// a not shared HTTP client with sensible timeouts.
package httpclient

import (
	"context"
	"net"
	"net/http"
	"time"
)

// DialContextFn was defined to make code more readable.
type dialContextFn func(ctx context.Context, network, address string) (net.Conn, error)

// DialContext implements our own dialer in order to set read and write idle timeouts.
func DialContext(rwtimeout, ctimeout time.Duration) dialContextFn {
	dialer := &net.Dialer{Timeout: ctimeout}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		c, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}

		if rwtimeout > 0 {
			timeoutConn := &tcpConn{
				TCPConn: c.(*net.TCPConn),
				timeout: rwtimeout,
			}
			return timeoutConn, nil
		}

		return c, nil
	}
}

// tcpConn is our own net.Conn which sets a read and write deadline and resets them each
// time there is read or write activity in the connection.
type tcpConn struct {
	*net.TCPConn
	timeout time.Duration
}

func (c *tcpConn) Read(b []byte) (int, error) {
	err := c.TCPConn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Read(b)
}

func (c *tcpConn) Write(b []byte) (int, error) {
	err := c.TCPConn.SetDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return 0, err
	}
	return c.TCPConn.Write(b)
}

// Default returns a default HTTP client with sensible values for slow 3G connections and above.
func Default() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext:           DialContext(30*time.Second, 10*time.Second),
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
}
