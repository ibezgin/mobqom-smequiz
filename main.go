package main

import (
	"flag"
	"fmt"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	fmt.Println("Hello world, smequiz!")
	http.ListenAndServe(*addr, nil)
}
