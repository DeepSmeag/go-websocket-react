package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		fmt.Println("Connected new on /echo:", conn.RemoteAddr())
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			fmt.Printf("%s sent %s\n", conn.RemoteAddr(), string(msg))
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
			if string(msg) == "exit" {
				conn.Close()
				break
			}
		}
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("Error:", err)
		}
		defer conn.Close()
		fmt.Println("Connected new on /ws:", conn.RemoteAddr())
		wg := sync.WaitGroup{}
		wg.Add(2)
		pauseChan := make(chan bool)
		go handleReading(conn, &wg, pauseChan)
		go handleWriting(conn, &wg, pauseChan)

		wg.Wait()

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}

func handleReading(conn *websocket.Conn, wg *sync.WaitGroup, pauseChan chan bool) {
	defer log.Println("Server stopped reader for", conn.RemoteAddr())
	defer wg.Done()
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message", err)
			return
		}
		cmd := string(msg)
		switch cmd {
		case "resume":
			pauseChan <- false
		case "pause":
			pauseChan <- true
		default:
			log.Println("Server received unknown command:", cmd)
		}
		log.Println("Server: received ", msgType, cmd)
	}
}
func handleWriting(conn *websocket.Conn, wg *sync.WaitGroup, pauseChan chan bool) {
	defer log.Println("Server stopped writer for", conn.RemoteAddr())
	defer wg.Done()
	paused := false
	for {
		select {
		case val := <-pauseChan:
			paused = val
		default:
			if !paused {
				randInt := strconv.Itoa(rand.IntN(10) + 1)
				if err := conn.WriteMessage(1, []byte(randInt)); err != nil {
					log.Println("ERROR WITH WRITING MESSAGE", err)
					return
				}
				if randInt == "10" {
					conn.Close()
					return
				}
			}
			time.Sleep(time.Second)
		}
	}
}
