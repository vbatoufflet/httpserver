# httpserver: HTTP server [![GoDoc][godoc-badge]][godoc-url]

Basic HTTP server for Go.

This server supports:

* TCP sockets
* UNIX sockets (with `user`, `group` and `mode` options)
* Graceful shutdown

## Example

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vbatoufflet/httpserver"
)

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hi there!\n"))
	})

	s := &httpserver.Server{
		Addr:            "localhost:8080", // or "unix:/path/to/server.sock?mode=0600"
		Handler:         m,
		ShutdownTimeout: 10 * time.Second,
	}

	// Wait for termination signal, then cancel context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	go func() { <-sigCh; cancel() }()

	// Start serving HTTP connections
	err := s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
```

[godoc-badge]: https://godoc.org/github.com/vbatoufflet/httpserver?status.svg
[godoc-url]: https://godoc.org/github.com/vbatoufflet/httpserver
