package client

import (
    "bufio"
    "net"
    "sync"
)

type Client struct {
    conn net.Conn
    reader *bufio.Reader
    mu sync.Mutex
}

func New_client(addr string) (*Client, error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil, err
    }

    return &Client{
        conn: conn,
        reader: bufio.NewReader(conn),
    }, nil
}

func (c *Client) Close() error {
    return c.conn.Close()
}
