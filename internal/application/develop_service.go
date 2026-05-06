package application

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type ListenerProvider interface {
	Listen(basePort int) (port int, ln net.Listener, err error)
}

type ServerStarter interface {
	Start(register func(mux *http.ServeMux), ln net.Listener) *http.Server
}

type AdminHandlerProvider interface {
	Handler() http.Handler
}

type APIRoutesProvider interface {
	Register(mux *http.ServeMux) error
}

type DevelopService struct {
	listener ListenerProvider
	server   ServerStarter
	admin    AdminHandlerProvider
	api      APIRoutesProvider
}

func NewDevelopService(listener ListenerProvider, server ServerStarter, admin AdminHandlerProvider, api APIRoutesProvider) *DevelopService {
	return &DevelopService{listener: listener, server: server, admin: admin, api: api}
}

func (s *DevelopService) Start(basePort int, host string) (url string, shutdown func(ctx context.Context) error, err error) {
	port, ln, err := s.listener.Listen(basePort)
	if err != nil {
		return "", nil, err
	}

	var registerErr error
	srv := s.server.Start(func(mux *http.ServeMux) {
		if s.api != nil {
			registerErr = s.api.Register(mux)
			if registerErr != nil {
				return
			}
		}
		mux.Handle("/", s.admin.Handler())
	}, ln)
	if registerErr != nil {
		_ = srv.Shutdown(context.Background())
		return "", nil, registerErr
	}

	url = fmt.Sprintf("http://%s:%d", host, port)
	shutdown = srv.Shutdown
	return url, shutdown, nil
}
