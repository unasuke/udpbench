package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type BenchResult struct {
	TransferreddByteSize int
	Count                int
	Failed               int
	TotalTime            time.Duration
}

func main() {
	var wait sync.WaitGroup
	var res BenchResult
	c := make(chan BenchResult)
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go func() {
			_, _ = sendrecv(10, c)
		}()
	}

	for {
		r := <-c
		res.Count += r.Count
		res.TotalTime += r.TotalTime
		res.Failed += r.Failed
		res.TransferreddByteSize += r.TransferreddByteSize
		if res.Count >= 100 {
			break
		}
	}
	fmt.Println(res)
}

func sendrecv(count int, ch chan BenchResult) (BenchResult, error) {
	var result BenchResult
	conn, err := net.Dial("udp", "0.0.0.0:8080")
	if err != nil {
		return result, err
	}

	for i := 0; i < count; i++ {
		uuid := uuid.NewString()
		s := time.Now()

		_, err := conn.Write([]byte(uuid))
		if err != nil {
			return result, err
		}

		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			return result, err
		}
		e := time.Now()

		result.Count++
		result.TotalTime += e.Sub(s)

		// Verify response
		if string(buf[:n]) != uuid {
			result.Failed++
		}
	}

	ch <- result

	return result, nil
}
