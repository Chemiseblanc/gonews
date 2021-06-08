package nntp

import (
	"crypto/tls"
	"log"
	"net"
)

type Server struct {
	Addr      string
	TLSConfig *tls.Config
	ErrorLog  *log.Logger
}

func NewServer(addr string, config *tls.Config) (Server, error) {
	srv := Server{
		Addr:      addr,
		TLSConfig: config,
	}
	return srv, nil
}

func (srv *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	for {
		rw, err := ln.Accept()
		if err != nil && srv.ErrorLog != nil {
			srv.ErrorLog.Println(err)
			continue
		}
		c := srv.newConn(rw)
		go c.serve()
	}
}

func (srv *Server) newConn(c net.Conn) Conn {
	return Conn{
		c,
		srv,
		State{},
	}
}
