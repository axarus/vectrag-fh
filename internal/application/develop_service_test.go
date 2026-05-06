package application

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"
)

// Mock implementations for testing
type mockListenerProvider struct {
	port int
	ln   net.Listener
	err  error
}

func (m *mockListenerProvider) Listen(basePort int) (port int, ln net.Listener, err error) {
	return m.port, m.ln, m.err
}

type mockServerStarter struct {
	server *http.Server
	err    error
}

func (m *mockServerStarter) Start(register func(mux *http.ServeMux), ln net.Listener) *http.Server {
	if m.err != nil {
		return nil
	}

	mux := http.NewServeMux()
	register(mux)

	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: mux,
	}

	go func() {
		server.Serve(ln)
	}()

	return server
}

type mockAdminHandlerProvider struct {
	handler http.Handler
}

func (m *mockAdminHandlerProvider) Handler() http.Handler {
	if m.handler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	}
	return m.handler
}

type mockAPIRoutesProvider struct {
	err error
}

func (m *mockAPIRoutesProvider) Register(mux *http.ServeMux) error {
	if m == nil {
		return nil
	}
	if m.err != nil {
		return m.err
	}

	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return nil
}

func TestNewDevelopService(t *testing.T) {
	listener := &mockListenerProvider{}
	server := &mockServerStarter{}
	admin := &mockAdminHandlerProvider{}
	api := &mockAPIRoutesProvider{}

	service := NewDevelopService(listener, server, admin, api)

	if service == nil {
		t.Fatal("NewDevelopService() returned nil")
	}

	if service.listener != listener {
		t.Error("listener not set correctly")
	}
	if service.server != server {
		t.Error("server not set correctly")
	}
	if service.admin != admin {
		t.Error("admin not set correctly")
	}
	if service.api != api {
		t.Error("api not set correctly")
	}
}

func TestDevelopService_Start_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer ln.Close()

	listener := &mockListenerProvider{
		port: 8080,
		ln:   ln,
		err:  nil,
	}

	server := &mockServerStarter{}
	admin := &mockAdminHandlerProvider{}
	api := &mockAPIRoutesProvider{}

	service := NewDevelopService(listener, server, admin, api)

	url, shutdown, err := service.Start(8080, "localhost")
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if url != "http://localhost:8080" {
		t.Errorf("Expected url = http://localhost:8080, got %s", url)
	}

	if shutdown == nil {
		t.Error("Expected shutdown function, got nil")
	}

	// Test shutdown function
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

func TestDevelopService_Start_ListenerError(t *testing.T) {
	listener := &mockListenerProvider{
		err: errors.New("listener error"),
	}

	server := &mockServerStarter{}
	admin := &mockAdminHandlerProvider{}
	api := &mockAPIRoutesProvider{}

	service := NewDevelopService(listener, server, admin, api)

	_, _, err := service.Start(8080, "localhost")
	if err == nil {
		t.Error("Expected error from listener, got nil")
	}

	if err.Error() != "listener error" {
		t.Errorf("Expected 'listener error', got %v", err)
	}
}

func TestDevelopService_Start_APIRegistrationError(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer ln.Close()

	listener := &mockListenerProvider{
		port: 8080,
		ln:   ln,
		err:  nil,
	}

	server := &mockServerStarter{}
	admin := &mockAdminHandlerProvider{}
	api := &mockAPIRoutesProvider{
		err: errors.New("api registration error"),
	}

	service := NewDevelopService(listener, server, admin, api)

	_, _, err = service.Start(8080, "localhost")
	if err == nil {
		t.Error("Expected error from API registration, got nil")
	}

	if err.Error() != "api registration error" {
		t.Errorf("Expected 'api registration error', got %v", err)
	}
}

func TestDevelopService_Start_NilAPI(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer ln.Close()

	listener := &mockListenerProvider{
		port: 8080,
		ln:   ln,
		err:  nil,
	}

	server := &mockServerStarter{}
	admin := &mockAdminHandlerProvider{}
	var api *mockAPIRoutesProvider = nil

	service := NewDevelopService(listener, server, admin, api)

	url, shutdown, err := service.Start(8080, "localhost")
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if url != "http://localhost:8080" {
		t.Errorf("Expected url = http://localhost:8080, got %s", url)
	}

	if shutdown == nil {
		t.Error("Expected shutdown function, got nil")
	}

	// Test shutdown function
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

func TestDevelopService_Start_DifferentHost(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer ln.Close()

	listener := &mockListenerProvider{
		port: 9090,
		ln:   ln,
		err:  nil,
	}

	server := &mockServerStarter{}
	admin := &mockAdminHandlerProvider{}
	api := &mockAPIRoutesProvider{}

	service := NewDevelopService(listener, server, admin, api)

	url, _, err := service.Start(9090, "example.com")
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	expectedURL := "http://example.com:9090"
	if url != expectedURL {
		t.Errorf("Expected url = %s, got %s", expectedURL, url)
	}
}
