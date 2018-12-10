package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Server is a HTTP server.
type Server struct {
	// Addr is the address for the server to listen on.
	// e.g. "localhost:8080" or "unix:/path/to/server.sock?mode=0600&user=www-data&group=www-data"
	Addr string

	// Handler is a HTTP handler serving requests.
	Handler http.Handler

	// ShutdownTimeout is the duration to wait before forcefully shuting down the server.
	ShutdownTimeout time.Duration
}

// Run starts accepting connections and serving HTTP requests.
func (s *Server) Run(ctx context.Context) error {
	socket, err := newSocket(s.Addr)
	if err != nil {
		return err
	}

	listener, err := net.Listen(socket.Proto, socket.Addr)
	if err != nil {
		return fmt.Errorf("cannot listen: %s", err)
	}
	defer listener.Close()

	err = socket.init()
	if err != nil {
		return fmt.Errorf("cannot initialize socket: %s", err)
	}

	server := &http.Server{
		Addr:    socket.Addr,
		Handler: s.Handler,
	}

	errCh := make(chan error)
	go func() {
		errCh <- server.Serve(listener)
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		// Stop immediately if no shutdown timeout
		if s.ShutdownTimeout == 0 {
			server.Shutdown(ctx)
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer cancel()

		server.Shutdown(ctx)
	}

	return nil
}
