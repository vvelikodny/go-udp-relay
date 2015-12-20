//
// Client to generate stream to test UDP relay
//

package main

import (
    "fmt"
    "net"
    "time"
    "os"
    "crypto/rand"
    "flag"
)

const (
    BUFF_SIZE = 512

    DFLT_SERVER_HOST = "localhost"
    DFLT_SERVER_PORT = 6666
)

func CheckError(err error) {
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(0)
    }
}

func randString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

    var bytes = make([]byte, n)

    rand.Read(bytes)

    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }

    return string(bytes)
}

func main() {
    var (
        serverHost = flag.String("server-host", DFLT_SERVER_HOST, "Outgoing stream host")
        serverPort = flag.Int("server-port", DFLT_SERVER_PORT, "Outgoing stream port")
        help = flag.Bool("h", false, "Show help")
    )

    flag.Parse()

    if *help {
        flag.Usage()
        os.Exit(0)
    }

    // Connect to server
    ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", *serverHost, *serverPort))
    CheckError(err)

    Conn, err := net.DialUDP("udp", nil, ServerAddr)
    CheckError(err)

    defer Conn.Close()

    buf := []byte(randString(BUFF_SIZE))

    packages := 0
    lastTime := time.Now()

    // Generate stream
    for {
        _, err := Conn.Write(buf)
        if err != nil {
            fmt.Println(err)
            continue;
        }

        if (time.Since(lastTime) >= time.Second) {
            fmt.Println(packages, " pps")

            lastTime = time.Now()
            packages = 0
        }

        packages++

        // Just to test pps
        time.Sleep(time.Nanosecond * 50000)
    }
}
