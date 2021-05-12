package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

const (
	MAX_PORT   = 1 << 16
	MAX_WORKER = 100
)

var IP = "localhost"
var wg sync.WaitGroup

func worker(ports <-chan int, result chan<- int) {
	defer wg.Done()
	for port := range ports {
		address := fmt.Sprintf("%s:%d", IP, port)
		if _, err := net.Dial("tcp", address); err == nil {
			result <- port
		} else {
			result <- -1
		}
	}
}

func main() {
	ports := make(chan int, MAX_WORKER)
	result := make(chan int)
	for i := 0; i < MAX_WORKER; i++ {
		wg.Add(1)
		go worker(ports, result)
	}
	go func() {
		for port := 0; port < MAX_PORT; port++ {
			ports <- port
		}
		close(ports)
	}()
	openPorts := []int{}
	for i := 0; i < MAX_PORT; i++ {
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
