package main

import (
	"log"
	"time"

	"github.com/openmind3d/rtpsocketgo"
	"github.com/pion/rtp"
)

func main() {
	sock, err := rtpsocketgo.Connect(rtpsocketgo.Config{
		Address:           "127.0.0.1:5500",
		UdpSocketMtuBytes: 1500,
	})

	if err != nil {
		log.Fatal(err)
	}

	rtpFanout := rtpsocketgo.NewFanoutFromSock(sock)

	rtpChan1 := make(chan *rtp.Packet)
	rtpChan2 := make(chan *rtp.Packet)

	rtpFanout.AddSub(rtpChan1)
	rtpFanout.AddSub(rtpChan2)

	go func() {
		for {
			_, ok := <-rtpChan1
			if !ok {
				log.Println("ch1 closed")
				return
			}
			log.Println("ch1")
		}
	}()

	go func() {
		for {
			_, ok := <-rtpChan2
			if !ok {
				log.Println("ch2 closed")
				return
			}
			log.Println("ch2")
		}
	}()

	time.Sleep(10 * time.Second)

	rtpFanout.Close()
	log.Println("fanout closed")

	time.Sleep(10 * time.Second)

	log.Println("exit")
}
