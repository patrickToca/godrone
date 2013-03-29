package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/felixge/godrone/src/attitude"
	"log"
	"net/http"
	"sync"
)

var clients []*websocket.Conn
var clientsLock sync.Mutex

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	go serveHttp()

	log.Printf("Initializing attitude ...")
	att, err := attitude.NewAttitude()
	if err != nil {
		panic(err)
	}

	log.Printf("Starting main loop ...")
	for {
		data, err := att.Update()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", data)
	}
}

func serveHttp() {
	http.Handle("/ws", websocket.Handler(handleWs))
	addr := ":80"
	log.Printf("serving clients at %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func handleWs(ws *websocket.Conn) {
	log.Printf("New client: %s", ws.RemoteAddr().String())
	clientsLock.Lock()
	clients = append(clients, ws)
	clientsLock.Unlock()

	var d string
	for {
		websocket.Message.Receive(ws, &d);
	}
}