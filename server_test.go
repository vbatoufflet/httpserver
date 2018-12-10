package httpserver

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Server_TCP(t *testing.T) {
	testServer(t,
		"localhost:8080",
		"tcp", "localhost:8080",
		nil,
	)
}

func Test_Server_UNIX(t *testing.T) {
	testServer(t,
		"unix:/tmp/httpserver-test.sock",
		"unix", "/tmp/httpserver-test.sock",
		nil,
	)
}

func Test_Server_UNIX_Params(t *testing.T) {
	testServer(
		t,
		"unix:/tmp/httpserver-test.sock?mode=0600",
		"unix", "/tmp/httpserver-test.sock",
		func(t *testing.T) {
			time.Sleep(20 * time.Millisecond)

			fi, err := os.Stat("/tmp/httpserver-test.sock")
			assert.Nil(t, err)
			assert.Equal(t, os.FileMode(0600)|os.ModeSocket, fi.Mode())
		},
	)
}

func testServer(t *testing.T, addr, netProto, netAddr string, test func(*testing.T)) {
	m := http.NewServeMux()
	m.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("foo")) })

	s := &Server{
		Addr:    addr,
		Handler: m,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if test != nil {
		go test(t)
	}

	go func() {
		time.Sleep(40 * time.Millisecond)

		client := &http.Client{
			Transport: &http.Transport{
				Dial: func(string, string) (net.Conn, error) {
					return net.Dial(netProto, netAddr)
				},
			},
		}

		resp, err := client.Get("http://foo/")
		assert.Nil(t, err)
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, []byte("foo"), b)

		cancel()
	}()

	err := s.Run(ctx)
	assert.Nil(t, err)
}
