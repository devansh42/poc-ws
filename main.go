package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {

	up := websocket.Upgrader{}
	up.ReadBufferSize = 1024
	up.WriteBufferSize = 1024
	up.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		conn, err := up.Upgrade(w, r, w.Header())
		if err != nil {
			log.Print(err)
			w.WriteHeader(500)
			return
		}
		mtype, reader, err := conn.NextReader()
		if err != nil {

		}
		switch mtype {
		case websocket.TextMessage:
			log.Print("Text Message")
			io.Copy(log.Writer(), reader)
		case websocket.BinaryMessage:
			log.Print("Binary Message")
			rr, _ := ioutil.ReadAll(reader)
			log.Print(string(rr))

		}
		go func(conn *websocket.Conn) {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
				wr, _ := conn.NextWriter(websocket.TextMessage)
				io.Copy(wr, strings.NewReader(fmt.Sprint("Message ", i)))
			}
		}(conn)

	})
	http.ListenAndServe(":8080", nil)

}
