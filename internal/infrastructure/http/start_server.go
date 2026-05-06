package http

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

type ServerStarter struct{}

func (ServerStarter) Start(register func(mux *http.ServeMux), ln net.Listener) *http.Server {
	return StartServer(register, ln)
}

func StartServer(register func(mux *http.ServeMux), ln net.Listener) *http.Server {

	mux := http.NewServeMux()
	if register != nil {
		register(mux)
	}

	server := &http.Server{
		Handler: mux,
	}

	// port := ln.Addr().(*net.TCPAddr).Port
	fmt.Println("Starting server")

	go func() {
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Error serving:", err)
		}
	}()

	return server
}
