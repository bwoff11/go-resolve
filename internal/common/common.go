package common

type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
	ProtocolDOT Protocol = "dot"
	ProtocolDOH Protocol = "doh"
)
