package server

import (
    "fmt"
    "net"
    "os"
)

func CheckError(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(1)
    }
}

const BUF_SIZE int = 512

type server struct {
    d *dispatcher
}

func New() server {
    return  server{}
}

func (s *server) Run(incomingPort int, outgoingPort int) {
    inputStreamAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", incomingPort))
    CheckError(err)

    outputStreamAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", outgoingPort))
    CheckError(err)

    // TODO refactor in & out channels on receiver & dispatcher
    s.d = NewDispatcher()

    // In-coming stream reader
    incomingStreamReceiver := NewReceiver(s.d.stream)
    go incomingStreamReceiver.Run(inputStreamAddr)

    // Out-coming stream dispatcher
    s.d.Run(outputStreamAddr)
}
