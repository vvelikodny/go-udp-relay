package server

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "net"
    "time"
)

var _ = Describe("Subscriber", func() {
    var (
        addr *net.UDPAddr
        s    *Subscriber
    )

    BeforeEach(func() {
        addr, _ = net.ResolveUDPAddr("udp", ":0")

        s = newSubscriber(addr)
    })

    It("Should have proper new state", func() {
        Expect(s.Address).Should(Equal(addr))

        Expect(s.lastActiveTime).Should(BeTemporally("~", time.Now(), time.Second))

        Expect(s.alive).ShouldNot(BeClosed())
        Expect(s.alive).ShouldNot(Receive())
    })

    It("Should update last alive time", func() {
        time.Sleep(time.Second)

        go s.UpdateAliveTime()

        Eventually(s.alive).Should(Receive(BeTemporally("~", time.Now(), time.Second)))
        Expect(s.alive).ShouldNot(BeClosed())
    })
})
