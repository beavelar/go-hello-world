package client

import (
	"fmt"
	"log"
	"simple_project/server"
	"time"

	"github.com/gorilla/websocket"
)

func receiveWsMessage(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("websocket client: error occurred reading websocket messages", err)
			return
		}
		log.Printf("websocket client: message received - %s\n", message)
	}
}

func StartWebsocketClient() {
	log.Println("websocket client: setting up client..")
	if server.WebsocketServerConfig == nil {
		log.Println("no websocket server config found, client won't be started..")
		return
	}

	path := fmt.Sprintf("%s://%s:%d/ws", server.WebsocketServerConfig.Protocol, server.WebsocketServerConfig.Host, server.WebsocketServerConfig.Port)
	conn, _, err := websocket.DefaultDialer.Dial(path, nil)

	if err != nil {
		log.Println("websocket client: failed to connect to websocket server", err)
		return
	}

	defer conn.Close()
	go receiveWsMessage(conn)

	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))

		if err != nil {
			log.Println("websocket client: error occurred writing message to websocket - ", err)
		}

		time.Sleep(time.Second)
	}
}
