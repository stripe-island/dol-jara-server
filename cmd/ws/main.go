package main

import (
	"flag"
	"log"
	"net/http"

	doljara "github.com/stripe-island/dol-jara-server"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	go doljara.HandleMessage()

	http.HandleFunc("/DoljaraRealtimeSync", doljara.ServWs)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
