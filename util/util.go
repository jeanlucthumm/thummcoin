package util

import (
	"net"
	"strconv"
)

func IPString(ip net.IP, port int) string {
	return ip.String() + ":" + strconv.Itoa(port)
}
