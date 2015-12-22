//
// Client to test UDP relay
//

package client

import (
    "errors"
    "fmt"
    "net"
    "time"
)

const (
    BUFF_SIZE = 512

    ALIVE_CHECK_TIME = time.Second * 30

    DFLT_SERVER_HOST = "localhost"
    DFLT_SERVER_PORT = 7777
)

type client struct {
    id   string
    Conn *net.UDPConn
}

func NewClient(id string) *client {
    return &client{id: id}
}

func (c *client) Start(localAddr *net.UDPAddr, serverAddr *net.UDPAddr) error {
    conn, err := net.DialUDP("udp", localAddr, serverAddr)
    if err != nil {
        return  err
    }

    c.Conn = conn

    return nil
}

func (c *client) Connect() error {
    return c.Send([]byte(fmt.Sprintf("CONNECT %v", c.id)))
}

func (c *client) Alive() error {
    return c.Send([]byte(fmt.Sprintf("ALIVE %v", c.id)))
}

func (c *client) Disconnect() error {
    if c.Conn == nil {
        return errors.New("Start before Disconnect!")
    }

    defer c.Conn.Close()

    err := c.Send([]byte(fmt.Sprintf("DISCONNECT %v", c.id)))
    if err != nil {
        return err
    }

    return nil
}

func (c *client) Send(data []byte) error {
    if c.Conn == nil {
        return errors.New("Start before send command!")
    }

    _, err := c.Conn.Write(data)
    if err != nil {
        return err
    }

    return nil
}
