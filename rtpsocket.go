package rtpsocketgo

import (
	"net"

	"github.com/pion/rtp"
)

const (
	DEFAULT_UDP_SOCKET_MTU_BYTES uint = 1500
)

type Config struct {
	Address           string
	UdpSocketMtuBytes uint
}

type RtpSocket struct {
	address string
	buf     []byte
	socket  *net.UDPConn
}

func Connect(config Config) (*RtpSocket, error) {
	udpSocketMtuBytes := DEFAULT_UDP_SOCKET_MTU_BYTES
	if config.UdpSocketMtuBytes != 0 {
		udpSocketMtuBytes = config.UdpSocketMtuBytes
	}

	s := &RtpSocket{
		address: config.Address,
		buf:     make([]byte, udpSocketMtuBytes),
	}

	if err := s.connectTo(config.Address); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *RtpSocket) connectTo(address string) error {
	listenAddress, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	sock, err := net.ListenUDP("udp", listenAddress)
	if err != nil {
		return err
	}
	s.socket = sock
	s.address = address

	return nil
}

func (s *RtpSocket) ReadRtpPacket() (*rtp.Packet, error) {
	n, _, err := s.socket.ReadFromUDP(s.buf)
	if err != nil {
		return nil, err
	}

	rtpPacket := rtp.Packet{}

	if err := rtpPacket.Unmarshal(s.buf[:n]); err != nil {
		return nil, err
	}

	return &rtpPacket, nil
}

func (s *RtpSocket) Close() error {
	return s.socket.Close()
}
