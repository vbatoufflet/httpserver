package httpserver

import (
	"net"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
)

type socket struct {
	Proto string
	Addr  string
	User  string
	Group string
	Mode  uint64
}

func newSocket(addr string) (*socket, error) {
	s := &socket{
		Proto: "tcp",
		Addr:  addr,
	}

	if strings.HasPrefix(s.Addr, "unix:") {
		s.Proto = "unix"
		s.Addr = strings.TrimPrefix(s.Addr, "unix:")

		if s.Addr == "" {
			return nil, ErrInvalidAddr
		}

		idx := strings.IndexByte(s.Addr, '?')
		if idx != -1 {
			values, err := url.ParseQuery(s.Addr[idx+1:])
			if err != nil {
				return nil, ErrInvalidAddr
			}

			s.Addr = s.Addr[:idx]
			s.User = values.Get("user")
			s.Group = values.Get("group")

			s.Mode, err = strconv.ParseUint(values.Get("mode"), 8, 32)
			if err != nil {
				return nil, ErrInvalidAddr
			}
		}
	} else {
		_, err := net.ResolveTCPAddr("tcp", s.Addr)
		if err != nil {
			return nil, ErrInvalidAddr
		}
	}

	return s, nil
}

func (s *socket) init() error {
	if s.Proto != "unix" {
		return nil
	}

	uid := os.Getuid()
	if s.User != "" {
		user, err := user.Lookup(s.User)
		if err != nil {
			return err
		}
		uid, _ = strconv.Atoi(user.Uid)
	}

	gid := os.Getgid()
	if s.Group != "" {
		group, err := user.LookupGroup(s.Group)
		if err != nil {
			return err
		}
		gid, _ = strconv.Atoi(group.Gid)
	}

	if err := os.Chown(s.Addr, uid, gid); err != nil {
		return err
	}

	if s.Mode != 0 {
		err := os.Chmod(s.Addr, os.FileMode(s.Mode))
		if err != nil {
			return err
		}
	}

	return nil
}
