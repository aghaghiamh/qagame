package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/gommon/log"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		defer conn.Close()

		readDone := make(chan bool)
		go ReadMessage(conn, readDone)
		
		writeDone := make(chan bool)
		msgCh := make(chan string)
		go produceMsg(msgCh, r.RemoteAddr)
		go WriteMessage(conn, msgCh, writeDone)

		<- readDone
	}))
}

func produceMsg(msgCh chan<- string, remoteAddr string) {
	for {
		msgCh <- remoteAddr
		time.Sleep(3*time.Second)
	}
}

func WriteMessage(conn net.Conn, msgCh <-chan string, writeDone chan<- bool) {
	for msg := range(msgCh) {
		err := wsutil.WriteServerMessage(conn, ws.OpText, []byte(msg))
		if err != nil {
			log.Error("err inside the readMessage: ", err)
			writeDone <- true
			return
		}
	}
}

func ReadMessage(conn net.Conn, done chan<- bool) {
	// enter following lines into command line: 
	// websocat ws://localhost:8080
	// {"event_type": "matched_palyers", "payload": "Hello to money!"}    
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			// handle error
			log.Error("err inside the readMessage: ", err)
			done <- true
			return
		}

		notif := &dto.Notification{}
		if err := json.Unmarshal(msg, notif); err != nil {
			log.Errorf("not being abled to unmarshal message with %s content", string(msg))
			done <- true
			return
		}

		fmt.Printf("notified %v with %d opCode\n", notif, op)
	}
}