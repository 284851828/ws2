package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func main() {
	ncount := 10000
	for i := 0; i < ncount; i++ {
		go newClient()
	}
	for true {
		time.Sleep(time.Second)
	}
}

func newClient() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/leffss"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go timeWriter(conn)

	for {
		_, _, err := conn.ReadMessage() // message
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		// fmt.Printf("received: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 5)
		// fmt.Printf("sended: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}
