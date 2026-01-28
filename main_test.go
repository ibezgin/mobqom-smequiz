package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const (
	HOST = "localhost"
)

type TestConfig struct {
	clientCount int
	wg          *sync.WaitGroup
}

func DialServer(wg *sync.WaitGroup) {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(fmt.Sprintf("ws://%s%s", HOST, WSPort), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		conn.Close()
		wg.Done()
	}()
	fmt.Println("connected to the server", conn.LocalAddr().String())
	time.Sleep(2 * time.Second)

}
func TestConnection(t *testing.T) {
	go createWSServer()
	time.Sleep(1 * time.Second)
	tc := TestConfig{
		clientCount: 5,
		wg:          new(sync.WaitGroup),
	}
	tc.wg.Add(tc.clientCount)

	for range tc.clientCount {
		go DialServer(tc.wg)
	}
	tc.wg.Wait()
	fmt.Println("exiting tests...")
}
