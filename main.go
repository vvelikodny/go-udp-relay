package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "github.com/vvelikodny/udprelay/server"
)

func CheckError(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(1)
    }
}

func printConnInfo(inputStreamPort *int, outputStreamPort *int) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Incoming stream can be send at the following addresses:")

    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                fmt.Printf("%v:%v\n", ipnet.IP, *inputStreamPort)
            }
        }
    }

    fmt.Println("Clients canconnect at the following addresses:")

    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                fmt.Printf("%v:%v\n", ipnet.IP, *outputStreamPort)
            }
        }
    }
}

func main() {
    var (
        help = flag.Bool("h", false, "Show help")
        incomingStreamPort = flag.Int("incoming-port", 6666, "Incoming stream port")
        outgoingStreamPort = flag.Int("outgoing-port", 7777, "Outgoing stream port")
    )

    flag.Parse()

    if *help {
        flag.Usage()
        os.Exit(0)
    }

    printConnInfo(incomingStreamPort, outgoingStreamPort)

    server := &server.Server{}
    server.Run(*incomingStreamPort, *outgoingStreamPort)

    os.Exit(0)
}
