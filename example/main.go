package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmcloughlin/profile"
)

func main() {
	fmt.Println("example UDP echo server")
	p := profile.Start()
	time.AfterFunc(30*time.Second, p.Stop)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGTERM,
		syscall.SIGINT)

	go func() {

		s := <-sig
		fmt.Println(s)
		os.Exit(0)
	}()

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
