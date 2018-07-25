package core

const (
	AddrTypeIPv4   = 0x01
	AddrTypeDomain = 0x03
	AddrTypeIPv6   = 0x04
)

type Host struct {
	Addr string
	Port string
	Type int32
}

func (h *Host) String() string {
	return h.Addr + ":" + h.Port
}
