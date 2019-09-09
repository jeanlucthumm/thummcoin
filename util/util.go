package util

import (
	"net"
)

func IPEqual(a *net.IPAddr, b *net.IPAddr) bool {
	return a.IP.Equal(b.IP) && a.Zone == b.Zone
}

func IPFromTCP(tcp *net.TCPAddr) *net.IPAddr {
	return &net.IPAddr{
		IP:   tcp.IP,
		Zone: tcp.Zone,
	}
}

func TCPFromIP(ip *net.IPAddr, port int) *net.TCPAddr {
	return &net.TCPAddr{
		IP:   ip.IP,
		Port: port,
		Zone: ip.Zone,
	}
}

func IPDialString(ip *net.IPAddr, port int) string {
	tcp := TCPFromIP(ip, port)
	return tcp.String()
}
