package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	flag "github.com/spf13/pflag"
)

type BenchResult struct {
	TransferreddByteSize int
	Count                int
	Failed               int
	TotalTime            time.Duration
}

type Option struct {
	Parallelism int
	Count       int
	Address     string
	Port        int
}

func main() {
	var res BenchResult
	var opt Option

	getOption(&opt)
	c := make(chan BenchResult)

	for i := 0; i < opt.Parallelism; i++ {
		go func() {
			_, _ = sendrecv(opt.Address+":"+strconv.Itoa(opt.Port), opt.Count, c)
		}()
	}

	for {
		r := <-c
		res.Count += r.Count
		res.TotalTime += r.TotalTime
		res.Failed += r.Failed
		res.TransferreddByteSize += r.TransferreddByteSize
		if res.Count >= opt.Parallelism*opt.Count {
			break
		}
	}

	printResult(res)
}

func getOption(opt *Option) {
	flag.IntVar(&opt.Parallelism, "parallelism", 10, "Worker parallelism number")
	flag.IntVar(&opt.Count, "count", 10, "Number of request from a worker")
	flag.StringVar(&opt.Address, "address", "127.0.0.1", "Server IP address")
	flag.IntVar(&opt.Port, "port", 8080, "Server port")
	flag.Parse()
}

func sendrecv(addr string, count int, ch chan BenchResult) (BenchResult, error) {
	var result BenchResult
	conn, err := net.Dial("udp", addr)
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

func printResult(opt BenchResult) {
	fmt.Printf("Total request count : %d\n", opt.Count)
	fmt.Printf("Total request time : %v\n", opt.TotalTime)
	fmt.Printf("Time per packets : %v\n", opt.TotalTime/time.Duration(opt.Count))
	fmt.Printf("Failed count : %d\n", opt.Failed)
}
