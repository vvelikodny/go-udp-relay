//
// Client to test UDP relay
//

package main

import (
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
    "flag"
    "time"
    "github.com/vvelikodny/udprelay/server"
)

func CheckError(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(0)
    }
}

const (
    BUFF_SIZE = 512

    ALIVE_CHECK_TIME = time.Second * 30

    DFLT_SERVER_HOST = "localhost"
    DFLT_SERVER_PORT = 7777
)

func main() {
    var (
        id = flag.String("id", "default", "id of the client")
        serverHost = flag.String("server-host", DFLT_SERVER_HOST, "Outgoing stream host")
        serverPort = flag.Int("server-port", DFLT_SERVER_PORT, "Outgoing stream port")
        help = flag.Bool("h", false, "Show help")
    )

    flag.Parse()

    if *help {
        flag.Usage()
        os.Exit(0)
    }

    LocalAddr, err := net.ResolveUDPAddr("udp", ":0")
    CheckError(err)

    ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", *serverHost, *serverPort))
    CheckError(err)

    client := server.NewClient(*id)
    client.Start(LocalAddr, ServerAddr)
    defer client.Disconnect()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)

    client.Connect()
    defer client.Disconnect()

    // Catch signals from os to exit
    go func() {
        <-c

        err = client.Disconnect()
        CheckError(err)

        client.Conn.Close()

        os.Exit(1)
    }()

    err = client.Connect()
    CheckError(err)

    // Alive status sender
    go func() {
        for {
            err := client.Alive()
            CheckError(err)

            time.Sleep(ALIVE_CHECK_TIME)
        }
    }()

    // Read stream from server
    buf := make([]byte, BUFF_SIZE)

    for {
        _, _, err := client.Conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Error: ", err)
            os.Exit(0)
        }

        fmt.Println("> " + string(buf[0:5]))
    }

    os.Exit(0)
}
