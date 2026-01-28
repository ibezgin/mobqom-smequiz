package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const (
	HOST = "localhost"
)

type TestConfig struct {
	clientCount    int
	wg             *sync.WaitGroup
	brMsgCount     atomic.Int64
	targetMsgCount int
}

func DialServer(tc *TestConfig) *websocket.Conn {
	exitCh := make(chan struct{})
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(fmt.Sprintf("ws://%s%s", HOST, WSPort), nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if int(tc.brMsgCount.Load()) == tc.targetMsgCount {
				close(exitCh)

				return
			}
		}

	}()

	go func() {
		<-exitCh
		conn.Close()
		tc.wg.Done()
	}()
	go func() {
		for {
			_, b, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if len(b) > 0 {
				tc.brMsgCount.Add(1)
			}

		}
	}()

	return conn
}
func TestConnection(t *testing.T) {
	go createWSServer()
	time.Sleep(1 * time.Second)
	clientCount := 5
	brCount := 3
	tc := &TestConfig{
		clientCount:    clientCount,
		wg:             new(sync.WaitGroup),
		brMsgCount:     atomic.Int64{},
		targetMsgCount: clientCount * brCount,
	}
	tc.wg.Add(tc.clientCount + 1)
	brClient := DialServer(tc)
	for range tc.clientCount {
		go DialServer(tc)
	}
	time.Sleep(1 * time.Second)

	for range brCount {
		msg := &ReqMsg{
			MsgType: MsgType_Broadcast,
			Data:    "hello from test",
		}
		time.Sleep(100 * time.Millisecond)

		err := brClient.WriteJSON(msg)
		if err != nil {
			fmt.Printf("error sending msg %v\n", err)
			return
		}
	}
	tc.wg.Wait()

	time.Sleep(1 * time.Second)

	fmt.Println("exiting tests...")
}
