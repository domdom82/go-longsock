package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"time"
)

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage: %s [server-url] \n", os.Args[0])
		fmt.Println("If server-url given, launch as client. Otherwise launch as server.")
		os.Exit(1)
	}

	flag.Parse()

	if flag.Arg(0) != "" {
		// client mode
		serverURL := flag.Arg(0)
		fmt.Printf("Connecting to %s\n", serverURL)

		conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
		if err != nil {
			panic(err)
		}
		handleConnection(conn)

	} else {
		//server mode
		fmt.Println("Starting server...")
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				panic(err)
			}
			handleConnection(conn)

		})
		http.ListenAndServe(":8080", mux)
	}
}

func handleConnection(conn *websocket.Conn) {
	defer func() {
		_ = conn.Close()
		fmt.Printf("Closed connection to %s\n", conn.RemoteAddr().String())
	}()

	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())

	for {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
			fmt.Println(err)
			return
		}
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("received '%s' from %s\n", p, conn.RemoteAddr().String())
		time.Sleep(time.Second)
	}

}
