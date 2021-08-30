package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("example UDP echo server")

	addr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8080,
	}

	serv, err := net.ListenUDP("udp", addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		buf := make([]byte, 512)
		n, ad, err := serv.ReadFromUDP(buf)

		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = serv.WriteToUDP(buf[:n], ad)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
