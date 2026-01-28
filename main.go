package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	WSPort = ":3223"
)

func createWSServer() {
	s := NewServer()
	go s.AcceptLoop()
	go fmt.Printf("Starting ws server on port %s\n", WSPort)
	http.HandleFunc("/", s.handleWs)
	log.Fatal(http.ListenAndServe(WSPort, nil))
}
func main() {
	createWSServer()
}
