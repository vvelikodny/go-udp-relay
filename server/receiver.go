package server

import (
    "net"
    "fmt"
)

type receiver struct {
    stream chan []byte
}

func NewReceiver(c chan []byte) *receiver {
    return &receiver{stream: c}
}

func (r *receiver) Run(addr *net.UDPAddr) {
    inputStreamConn, err := net.ListenUDP("udp", addr)
    CheckError(err)

    defer inputStreamConn.Close()

    buf := make([]byte, 1024)

    // Read incoming stream
    for {
        n, _, err := inputStreamConn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Error: ", err)
            continue
        }

        r.stream <- buf[0:n]
    }
}