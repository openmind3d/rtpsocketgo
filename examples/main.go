package main

import (
	"log"
	"time"

	rtpsocket "github.com/openmind3d/rtpsocketgo"
)

func main() {
	rtpSocket, err := rtpsocket.Connect(rtpsocket.Config{
		Address:           "127.0.0.1:5500",
		UdpSocketMtuBytes: 1500,
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			rtpPacket, err := rtpSocket.ReadRtpPacket()
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Have packet: ", rtpPacket.String())
		}
	}()

	<-time.After(15 * time.Second)
	rtpSocket.Close()

	select {}
}
