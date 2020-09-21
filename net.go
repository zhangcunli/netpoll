// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package netpoll

import (
	"errors"
	"io"
	"net"
	"syscall"
)

const (
	bufferSize = 0x10000
)

// EOF is the error returned by Read when no more input is available.
var EOF = io.EOF

// EAGAIN is the error when resource temporarily unavailable
var EAGAIN = syscall.EAGAIN

// ErrServerClosed is returned by the Server's Serve and ListenAndServe
// methods after a call to Close.
var ErrServerClosed = errors.New("Server closed")

// ListenAndServe listens on the network address and then calls
// Serve with handler to handle requests on incoming connections.
//
// The handler must be not nil.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(network, address string, handler Handler) error {
	server := &Server{Network: network, Address: address, Handler: handler}
	return server.ListenAndServe()
}

// Serve accepts incoming connections on the listener l,
// and registers the conn fd to poll. The poll will trigger the fd to
// read requests and then call handler to reply to them.
//
// The handler must be not nil.
//
// Serve always returns a non-nil error.
func Serve(lis net.Listener, handler Handler) error {
	server := &Server{Handler: handler}
	return server.Serve(lis)
}

type netServer struct {
	listener net.Listener
	Handler  Handler
}

func (s *netServer) Serve(l net.Listener) (err error) {
	s.listener = l
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer func() {
				if e := recover(); e != nil {
				}
			}()
			var err error
			var context Context
			if context, err = s.Handler.Upgrade(c); err != nil {
				c.Close()
				return
			}
			for err == nil {
				err = s.Handler.Serve(context)
			}
			c.Close()
		}(conn)
	}
	return nil
}

func (l *netServer) Close() error {
	return l.listener.Close()
}
