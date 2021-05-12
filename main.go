package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	MAX_WORKER = 100
)

var IP = "localhost"
var wg sync.WaitGroup
var timeout = 0

func worker(ports <-chan int, result chan<- int) {
	defer wg.Done()
	for port := range ports {
		address := fmt.Sprintf("%s:%d", IP, port)
		if _, err := net.DialTimeout("tcp", address, time.Millisecond*time.Duration(timeout)); err == nil {
			result <- port
		} else {
			result <- -1
		}
	}
}

func main() {
	start := 0
	end := 65535
	flag.IntVar(&start, "start", 0, "start")
	flag.IntVar(&end, "end", 65535, "end")
	flag.StringVar(&IP, "ip", "localhost", "IP")
	flag.IntVar(&timeout, "timeout", 3, "timeout(ms)")
	flag.Parse()

	if start < 0 || start > 65535 || end < 0 || end > 65535 {
		panic(errors.New("port range error"))
	}
	ports := make(chan int, MAX_WORKER)
	result := make(chan int)
	for i := 0; i < MAX_WORKER; i++ {
		wg.Add(1)
		go worker(ports, result)
	}
	go func() {
		for port := start; port <= end; port++ {
			ports <- port
		}
		close(ports)
	}()
	openPorts := []int{}
	for i := start; i <= end; i++ {
		port := <-result
		if port != -1 {
			openPorts = append(openPorts, port)
		}
	}
	sort.Ints(openPorts)
	fmt.Println(openPorts)
	wg.Wait()
}

//func main() {
//	wg := sync.WaitGroup{}
//	mu := sync.Mutex{}
//	openPorts := []int{}
//	for i := 0; i < MAX_PORT; i++ {
//		wg.Add(1)
//		go func(port int) {
//			defer wg.Done()
//			address := fmt.Sprintf("%s:%d", IP, port)
//			if _, err := net.Dial("tcp", address); err == nil {
//				mu.Lock()
//				openPorts = append(openPorts, port)
//				mu.Unlock()
//			}
//		}(i)
//	}
//	wg.Wait()
//	fmt.Println(openPorts)
//}
