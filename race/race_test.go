package race

import (
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
)

func TestRace(t *testing.T){
		CanConflict()
}

func CanConflict(){
	c := make(chan bool)
	m := make(map[string]string)
	go func() {
		m["1"] = "a" // First conflicting access.
		c <- true
	}()
	m["2"] = "b" // Second conflicting access.
	<-c
	for k, v := range m {
		fmt.Println(k, v)
	}
}

func LoopCounterRace(){
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i) // Not the 'i' you are looking for.
			wg.Done()
		}()
	}
	wg.Wait()
	/* Correction:
	go func(j int) {
				fmt.Println(j) // Good. Read local copy of the loop counter.
				wg.Done()
			}(i)
	*/
}

func ShareVariableRace(){
	data := []byte{43,141,43,141,41,44,11}
	res := make(chan error, 2)
	f1, err := os.Create("file1")
	if err != nil {
		res <- err
	} else {
		go func() {
			// This err is shared with the main goroutine,
			// so the write races with the write below.
			_, err = f1.Write(data)
			res <- err
			f1.Close()
		}()
	}
	f2, err := os.Create("file2") // The second conflicting write to err.
	if err != nil {
		res <- err
	} else {
		go func() {
			_, err = f2.Write(data)
			res <- err
			f2.Close()
		}()
	}
	/*
	fixed:
	_, err := f1.Write(data)
				...
	_, err := f2.Write(data)
				...
	*/
}

// unprotected map
// use mutex: serviceMu sync.Mutex
var service map[string]net.Addr

func RegisterService(name string, addr net.Addr) {
	service[name] = addr
}

func LookupService(name string) net.Addr {
	return service[name]
}

// unprotected primitive type can also cause a race which can be fixed by built-in atomic operations in sync/atomic package.

func UnsynchronizedSendReceive(){
	c := make(chan struct{}) // or buffered channel

	// The race detector cannot derive the happens before relation
	// for the following send and close operations. These two operations
	// are unsynchronized and happen concurrently.
	go func() { c <- struct{}{} }()
	// fixed by : <-c
	close(c)
}