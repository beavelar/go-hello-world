package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func httpReq(w http.ResponseWriter, r *http.Request) {
	log.Println("http server: /http request received")
	io.WriteString(w, "simple server http response")
}

func redirectedReq(w http.ResponseWriter, r *http.Request) {
	log.Println("http server: /redirected request received")
	io.WriteString(w, "simple server redirected response")
}

func rootReq(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		redirect := fmt.Sprintf("%s://%s:%d", HttpServerConfig.Protocol, HttpServerConfig.Host, HttpServerConfig.Port)
		http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
	} else {
		log.Println("http server: / request received")
		io.WriteString(w, "simple server root response")
	}
}

func StartHttpServer() {
	log.Println("http server: setting up server..")
	if HttpServerConfig == nil {
		HttpServerConfig = &Config{Host: "localhost", Port: 55571, Protocol: "http"}
	}

	mux := http.NewServeMux()
	redirect := fmt.Sprintf("%s://%s:%d/redirected", HttpServerConfig.Protocol, HttpServerConfig.Host, HttpServerConfig.Port)

	nh := http.NotFoundHandler()
	rh := http.RedirectHandler(redirect, http.StatusPermanentRedirect)

	mux.Handle("/redirect", rh)
	mux.Handle("/not-found", nh)

	mux.HandleFunc("/", rootReq)
	mux.HandleFunc("/http", httpReq)
	mux.HandleFunc("/redirected", redirectedReq)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", HttpServerConfig.Host, HttpServerConfig.Port),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       10 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("http server: received server closed error - %v\n", err)
	} else if err != nil {
		log.Printf("http server: error occurred on http http server - %v\n", err)
	}
}
