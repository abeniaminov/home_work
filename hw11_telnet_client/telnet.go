package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	con     net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address, timeout: timeout, in: in, out: out,
	}
}

func (c *client) Connect() error {
	var err error
	c.con, err = net.DialTimeout("tcp", c.address, c.timeout)
	return err
}

func (c *client) Send() error {
	_, err := io.Copy(c.con, c.in)
	return err
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.con)
	return err
}

func (c *client) Close() error {
	return c.con.Close()
}
