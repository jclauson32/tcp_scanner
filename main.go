package main

import (
	"fmt"
	"net"
	"sort"
)

/*
*Worker goroutine used to process port numbers

ports (chan int): input channel used to work through ports

results (chan int): output channel used to return open ports or errors (0)
*/
func worker(ports chan int, results chan int) {

	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println(err)
			results <- 0
			continue
		}

		conn.Close()
		results <- p
	}
}

/*
*
Main method that initializes input (ports) and output (results) channels and initializes worker from worker class.
Loops through all possible ports 1 - 65535 and sends them to the input channel.
Loops through all possible to make sure every port is acked.
Prints open ports in order
*/
func main() {
	ports := make(chan int, 1000)
	results := make(chan int)
	var openports []int
	for i := 1; i < cap(ports); i++ {
		go worker(ports, results)

	}
	go func() {
		for i := 1; i <= 65535; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 65535; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
