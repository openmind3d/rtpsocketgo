package rtpsocketgo

import (
	"errors"
	"net"

	"github.com/pion/rtp"
)

type RtpFanout struct {
	sock    *RtpSocket
	subs    map[chan *rtp.Packet]struct{}
	subCh   chan chan *rtp.Packet
	unsubCh chan chan *rtp.Packet
	closeCh chan struct{}
}

func NewFanoutFromSock(sock *RtpSocket) *RtpFanout {
	f := &RtpFanout{
		sock:    sock,
		subs:    map[chan *rtp.Packet]struct{}{},
		subCh:   make(chan chan *rtp.Packet, 10),
		unsubCh: make(chan chan *rtp.Packet, 10),
		closeCh: make(chan struct{}, 1),
	}

	inRtpPkt := make(chan *rtp.Packet, 10)

	go func() {
		for {
			rtpPacket, err := f.sock.ReadRtpPacket()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					close(inRtpPkt)
					f.closeCh <- struct{}{}
					return
				}
			}
			if rtpPacket != nil {
				inRtpPkt <- rtpPacket
			}
		}
	}()

	go func() {
		for {
			select {
			case rtpPacket, ok := <-inRtpPkt:
				if ok {
					for sub := range f.subs {
						select {
						case sub <- rtpPacket.Clone(): // deep copy packet!!!
						default: // skip if full
						}
					}
				}
			case sub, ok := <-f.subCh:
				if ok {
					f.subs[sub] = struct{}{}
				}
			case sub, ok := <-f.unsubCh:
				if ok {
					delete(f.subs, sub)
				}
			case <-f.closeCh:
				for sub := range f.subs {
					f.unsubCh <- sub
					close(sub)
				}
				return
			}
		}
	}()

	return f
}

func (f *RtpFanout) AddSub(sub chan *rtp.Packet) {
	f.subCh <- sub
}

func (f *RtpFanout) DelSub(sub chan *rtp.Packet) {
	f.unsubCh <- sub
}

func (f *RtpFanout) Close() {
	f.sock.Close() // will be handled correctly by reading goroutine
}
