package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "github.com/vvelikodny/udprelay/server"
)

func printConnInfo(inputStreamPort int, outputStreamPort int) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Incoming stream can be send at the following addresses:")
    printAddrPortInfo(addrs, inputStreamPort)

    fmt.Println("Clients canconnect at the following addresses:")
    printAddrPortInfo(addrs, outputStreamPort)
}

func printAddrPortInfo(addrs []net.Addr, port int) {
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                fmt.Printf("%v:%v\n", ipnet.IP, port)
            }
        }
    }
}

func main() {
    var help bool

    flag.BoolVar(&help, "h", false, "Show help")

    var incomingStreamPort int
    var outgoingStreamPort int

    flag.IntVar(&incomingStreamPort, "incoming-port", 6666, "Incoming stream port")
    flag.IntVar(&outgoingStreamPort, "outgoing-port", 7777, "Outgoing stream port")

    flag.Parse()

    if help {
        flag.Usage()
        os.Exit(0)
    }

    printConnInfo(incomingStreamPort, outgoingStreamPort)

    server := server.New()
    server.Run(incomingStreamPort, outgoingStreamPort)
}
