package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"simple_project/server"
	"time"
)

func httpRequest(path string) {
	if res, err := http.Get(fmt.Sprintf("%s/http", path)); err == nil {
		defer res.Body.Close()
		if resBytes, err := io.ReadAll(res.Body); err == nil {
			log.Printf("http client: /http received response from server - %s\n", string(resBytes))
		}
	}
}

func notFoundRequest(path string) {
	if res, err := http.Get(fmt.Sprintf("%s/not-found", path)); err == nil {
		defer res.Body.Close()
		if resBytes, err := io.ReadAll(res.Body); err == nil {
			log.Printf("http client: /not-found received response from server - %s\n", string(resBytes))
		}
	}
}

func redirectRequest(path string) {
	if res, err := http.Get(fmt.Sprintf("%s/redirect", path)); err == nil {
		defer res.Body.Close()
		if resBytes, err := io.ReadAll(res.Body); err == nil {
			log.Printf("http client: /redirect received response from server - %s\n", string(resBytes))
		}
	}
}

func rootRequest(path string) {
	if res, err := http.Get(path); err == nil {
		defer res.Body.Close()
		if resBytes, err := io.ReadAll(res.Body); err == nil {
			log.Printf("http client: / received response from server - %s\n", string(resBytes))
		}
	}
}

func StartHttpClient() {
	log.Println("http client: setting up client..")
	if server.HttpServerConfig == nil {
		log.Println("no server config found, client won't be started..")
		return
	}

	path := fmt.Sprintf("%s://%s:%d", server.HttpServerConfig.Protocol, server.HttpServerConfig.Host, server.HttpServerConfig.Port)
	for {
		go httpRequest(path)
		go notFoundRequest(path)
		go redirectRequest(path)
		go rootRequest(path)

		time.Sleep(time.Second)
	}
}
