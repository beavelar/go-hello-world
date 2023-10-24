package main

import (
	"simple_project/client"
	"simple_project/server"
	"time"
)

func main() {
	server.GrpcServerConfig = &server.Config{Host: "localhost", Port: 55571}
	server.HttpServerConfig = &server.Config{Host: "localhost", Port: 55572, Protocol: "http"}
	server.WebsocketServerConfig = &server.Config{Host: "localhost", Port: 55573, Protocol: "ws"}

	go server.StartHttpServer()
	go server.StartWebsocketServer()
	go server.StartGrpcServer()

	time.Sleep(2 * time.Second)

	go client.StartHttpClient()
	go client.StartWebsocketClient()
	go client.StartGrpcClient()

	select {}
}
