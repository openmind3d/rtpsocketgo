# Rtp socket
Receive rtp packets from udp socket and send to subscribers
## Installing
```sh
go get github.com/openmind3d/rtpsocketgo@latest
```
## Rtp socket
```go
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
```
## Rtp fanout
```go
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

```