package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Socket_TCP_Simple(t *testing.T) {
	s, err := newSocket("localhost:80")
	assert.Nil(t, err)
	assert.Equal(t, "localhost:80", s.Addr)
	assert.Equal(t, "", s.User)
	assert.Equal(t, "", s.Group)
	assert.Equal(t, uint64(0), s.Mode)
}

func Test_Socket_TCP_Invalid(t *testing.T) {
	_, err := newSocket("invalid")
	assert.Equal(t, ErrInvalidAddr, err)
}

func Test_Socket_UNIX_Simple(t *testing.T) {
	s, err := newSocket("unix:/path/to/server.sock")
	assert.Nil(t, err)
	assert.Equal(t, "/path/to/server.sock", s.Addr)
	assert.Equal(t, "", s.User)
	assert.Equal(t, "", s.Group)
	assert.Equal(t, uint64(0), s.Mode)
}

func Test_Socket_UNIX_Params(t *testing.T) {
	s, err := newSocket("unix:/path/to/server.sock?user=nobody&group=nogroup&mode=0600")
	assert.Nil(t, err)
	assert.Equal(t, "/path/to/server.sock", s.Addr)
	assert.Equal(t, "nobody", s.User)
	assert.Equal(t, "nogroup", s.Group)
	assert.Equal(t, uint64(0600), s.Mode)
}

func Test_Socket_UNIX_Empty(t *testing.T) {
	_, err := newSocket("unix:")
	assert.Equal(t, ErrInvalidAddr, err)
}

func Test_Socket_UNIX_InvalidMode(t *testing.T) {
	_, err := newSocket("unix:/path/to/server.sock?mode=a")
	assert.Equal(t, ErrInvalidAddr, err)
}
