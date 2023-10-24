package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	ReadBufferSize:   5,
	WriteBufferSize:  5,
	CheckOrigin:      checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func wsRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("websocket server: /ws request received")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("websocket server: upgrade failed", err)
		return
	}

	defer log.Println("websocket server: closing websocket connection")
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("websocket server: read failed", err)
			break
		}

		input := string(message)
		log.Printf("websocket server: received message - %s\n", input)

		message = []byte("pong")
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("websocket server: write failed", err)
			break
		}
	}
}

func StartWebsocketServer() {
	log.Println("websocket server: setting up server..")
	if WebsocketServerConfig == nil {
		WebsocketServerConfig = &Config{Host: "localhost", Port: 55572, Protocol: "ws"}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsRequest)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", WebsocketServerConfig.Host, WebsocketServerConfig.Port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       10 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("websocket server: received server closed error - %v\n", err)
	} else if err != nil {
		log.Printf("websocket server: error occurred on http server - %v\n", err)
	}
}
