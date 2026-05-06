package http

import (
	"fmt"
	"net"
)

type ListenerProvider struct{}

func (ListenerProvider) Listen(basePort int) (int, net.Listener, error) {
	return GetAvailablePort(basePort)
}

func GetAvailablePort(basePort int) (int, net.Listener, error) {
	var ln net.Listener
	var err error

	for i := 0; i < 3; i++ {
		tryPort := basePort + i
		ln, err = net.Listen("tcp", fmt.Sprintf(":%d", tryPort))
		if err == nil {
			return tryPort, ln, nil
		}
	}

	return 0, nil, fmt.Errorf("all three ports are in use; app cannot be started")
}
