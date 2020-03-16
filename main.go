package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	retry := flag.Bool("retry", false, "Reconnect after disconnection (client only)")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] [server-url] \n", os.Args[0])
		fmt.Println("If server-url given, launch as client. Otherwise launch as server.")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.LUTC)
	log.SetOutput(os.Stdout)

	if flag.Arg(0) != "" {
		// client mode
		serverURL := flag.Arg(0)
		log.Printf("Connecting to %s\n", serverURL)

		connect := func() error {
			conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
			if err != nil {
				log.Println(err)
				return err
			}
			return handleConnection(conn)
		}
		if *retry {
			withRetry(connect)
		} else {
			connect()
		}

	} else {
		//server mode
		log.Println("Starting server...")
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if forwardedFor := r.Header.Get("x-forwarded-for"); forwardedFor != "" {
				log.Printf("x-forwarded-for: %s\n", forwardedFor)
			}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				panic(err)
			}
			_ = handleConnection(conn)

		})
		_ = http.ListenAndServe(":8080", mux)
	}
}

func handleConnection(conn *websocket.Conn) error {
	defer func() {
		_ = conn.Close()
		log.Printf("Closed connection: %s -> %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())
	}()

	log.Printf("New connection: %s -> %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())

	for i := 0; ; i++ {
		msg := fmt.Sprintf("ping %d", i)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println(err)
			return err
		}
		log.Printf("-> Sent '%s' to %s\n", msg, conn.RemoteAddr().String())
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("<- Received '%s' from %s\n", p, conn.RemoteAddr().String())
		time.Sleep(time.Second)
	}

}

func withRetry(f func() error) {
	err := f()

	if err != nil {
		time.Sleep(time.Second)
		withRetry(f)
	}
}
