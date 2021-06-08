package nntp

import (
	"bufio"
	"crypto/tls"
	"errors"
	"net"
	"strings"
)

type Conn struct {
	rwc net.Conn

	server *Server
	state  State
}

var commandMap = map[string]func(*Conn, []string) error{
	"ARTICLE":      (*Conn).article,
	"BODY":         (*Conn).body,
	"CAPABILITIES": (*Conn).capabilities,
	"DATE":         (*Conn).date,
	"GROUP":        (*Conn).group,
	"HDR":          (*Conn).hdr,
	"HEAD":         (*Conn).head,
	"HELP":         (*Conn).help,
	"IHAVE":        (*Conn).ihave,
	"LAST":         (*Conn).last,
	"LIST":         (*Conn).list,
	"LISTGROUP":    (*Conn).listgroup,
	"MODE":         (*Conn).mode,
	"NEWGROUPS":    (*Conn).newgroups,
	"NEWNEWS":      (*Conn).newnews,
	"NEXT":         (*Conn).next,
	"OVER":         (*Conn).over,
	"POST":         (*Conn).post,
	"STAT":         (*Conn).stat,
	"QUIT":         (*Conn).quit,
	"STARTTLS":     (*Conn).starttls,
}

func require_arg_length(args []string, length int) error {
	if len(args) != length {
		return errors.New(command_syntax_error)
	} else {
		return nil
	}
}

func (conn *Conn) serve() {
	s := bufio.NewScanner(conn.rwc)

	for s.Scan() {
		cmd_args := strings.Fields(s.Text())
		cmd, args := cmd_args[0], cmd_args[1:]
		cmd = strings.ToUpper(cmd)
		if handler, ok := commandMap[cmd]; ok {
			err := handler(conn, args)
			if err != nil {
				conn.rwc.Write([]byte(err.Error()))
			}
		}
	}
}

func (c *Conn) group(args []string) error {
	return nil

}
func (c *Conn) article(args []string) error {
	return nil
}
func (c *Conn) head(args []string) error {
	return nil
}
func (c *Conn) body(args []string) error {
	return nil
}
func (c *Conn) stat(args []string) error {
	return nil
}
func (c *Conn) last(args []string) error {
	return nil
}
func (c *Conn) next(args []string) error {
	return nil
}
func (c *Conn) post(args []string) error {
	return nil
}
func (c *Conn) ihave(args []string) error {
	return nil
}
func (c *Conn) newgroups(args []string) error {
	return nil
}
func (c *Conn) newnews(args []string) error {
	return nil
}

func (c *Conn) quit(args []string) error {
	return nil
}
func (c *Conn) starttls(args []string) error {
	if c.server.TLSConfig != nil {
		tlsConn := tls.Server(c.rwc, c.server.TLSConfig)
		c.rwc = net.Conn(tlsConn)
		c.serve()
	}
	return nil
}
