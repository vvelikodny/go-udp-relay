package server

import (
    "testing"
    "fmt"
    "net"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "time"
)

const (
    IN_COMING_PORT = 6666
    OUT_COMING_PORT = 7777
)

func TestServerStart(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Server Suite")
}

var _ = Describe("Server", func() {
    server := &Server{}

    client1 := NewClient("client1")
    client2 := NewClient("client2")

    go server.Run(IN_COMING_PORT, OUT_COMING_PORT)

    BeforeEach(func() {

    })

    It("Should have zero subscribers on start", func() {
        Expect(server.d.subscribers).Should(HaveLen(0))
    })

    It("Should have zero subscribers on client socket connect", func() {
        LocalAddr, err := net.ResolveUDPAddr("udp", ":0")
        CheckError(err)

        ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", "localhost", OUT_COMING_PORT))
        CheckError(err)

        client1.Start(LocalAddr, ServerAddr)

        Expect(server.d.subscribers).Should(HaveLen(0))
    })

    It("Shouldn't add subscriber on client send ALIVE command before CONNECT", func() {
        client1.Alive()

        time.Sleep(time.Second)

        Expect(server.d.subscribers).Should(HaveLen(0))
    })

    It("Should have add subscriber on client send CONNECT command", func() {
        client1.Connect()

        time.Sleep(time.Second)

        Expect(server.d.subscribers).Should(HaveLen(1))
    })

    It("Should have one client then second just connect to socket", func() {
        LocalAddr, err := net.ResolveUDPAddr("udp", ":0")
        CheckError(err)

        ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", "localhost", OUT_COMING_PORT))
        CheckError(err)

        client2.Start(LocalAddr, ServerAddr)

        time.Sleep(time.Second)

        Expect(server.d.subscribers).Should(HaveLen(1))
    })

    It("Should have two subscribers then second client send CONNECT command", func() {
        client2.Connect()

        time.Sleep(time.Second)

        Expect(server.d.subscribers).Should(HaveLen(2))
    })

    It("Should have delete subscriber on client send DISCONNECT command", func() {
        client1.Disconnect()

        time.Sleep(time.Second)

        Expect(server.d.subscribers).Should(HaveLen(1))
    })

    It("Should have delete subscriber on client send DISCONNECT command", func() {
        client2.Disconnect()

        time.Sleep(time.Second)

        Expect(server.d.subscribers).Should(HaveLen(0))
    })
})