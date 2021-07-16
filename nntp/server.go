package nntp

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"strings"
	"time"
)

type Peer struct {
	name     string
	addr     net.Addr
	lastsync time.Time
	groups   []Group
}

type Server struct {
	Addr      string
	TLSConfig *tls.Config
	Log       *log.Logger

	storage Storage
	auth    Auth
	filter  FilterFunc

	peers []Peer
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
		if err != nil && srv.Log != nil {
			srv.Log.Println(err)
			continue
		}
		c := srv.NewConn(rw)
		go srv.serve(c)
	}
}

func (*Server) serve(c *Conn) {
	s := bufio.NewScanner(c.br)

	tlsStatus := c.isTLS

	for s.Scan() {
		cmd_args := strings.Fields(s.Text())
		cmd, args := cmd_args[0], cmd_args[1:]
		cmd = strings.ToUpper(cmd)
		if handler, ok := commandMap[cmd]; ok {
			if err := handler(c, args); err != nil {
				if e, ok := err.(net.Error); ok && !e.Temporary() {
					c.server.Log.Printf(e.Error())
					c.Close()
				}
			}

			// Renew the scanner if the connection is upgraded to use TLS
			if !tlsStatus && c.isTLS {
				s = bufio.NewScanner(c.br)
				tlsStatus = true
			}
		}
	}
}

func (srv *Server) NewConn(c net.Conn) *Conn {
	_, isTLS := c.(*tls.Conn)

	return &Conn{
		c,
		bufio.NewReader(c),
		bufio.NewWriter(c),

		isTLS,

		srv,
		nil,
		nil,
	}
}

func (srv *Server) PropagateNews() {

}
