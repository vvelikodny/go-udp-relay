package server

import (
    "net"
    "time"
)

const (
    TTL = time.Second * 30
)

type Subscriber struct {
    Address        *net.UDPAddr
    lastActiveTime time.Time
    alive          chan time.Time
}

func newSubscriber(address *net.UDPAddr) *Subscriber {
    return &Subscriber{Address: address, lastActiveTime: time.Now(), alive: make(chan time.Time)}
}

func (s *Subscriber) UpdateAliveTime() {
    s.alive <- time.Now();
}

func (s *Subscriber) CheckAlive(unsubscribe chan *net.UDPAddr) {
    aliveCheckTicker := time.NewTicker(TTL)

    for {
        select {
        case <-aliveCheckTicker.C:
            if time.Since(s.lastActiveTime) >= TTL {
                aliveCheckTicker.Stop()
                unsubscribe <- s.Address
                return
            }
        case t := <- s.alive:
            s.lastActiveTime = t
            aliveCheckTicker.Stop()
            aliveCheckTicker = time.NewTicker(TTL)
        }
    }
}