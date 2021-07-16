package nntp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/textproto"
)

// Conn is a stateful connection that that allows for buffered IO
type Conn struct {
	net.Conn
	br *bufio.Reader
	bw *bufio.Writer

	isTLS bool

	server        *Server
	articleNumber *uint
	group         *Group
}

var commandMap = map[string]func(*Conn, []string) error{
	"ARTICLE":      ArticleHander,
	"BODY":         BodyHandler,
	"CAPABILITIES": CapabilitiesHandler,
	"DATE":         DateHandler,
	"GROUP":        GroupHandler,
	"HDR":          HdrHandler,
	"HEAD":         HeadHandler,
	"HELP":         HelpHandler,
	"IHAVE":        IhaveHandler,
	"LAST":         LastHandler,
	"LIST":         ListHandler,
	"LISTGROUP":    ListgroupHandler,
	"MODE":         ModeHandler,
	"NEWGROUPS":    NewgroupsHandler,
	"NEWNEWS":      NewnewsHandler,
	"NEXT":         NextHandler,
	"OVER":         OverHandler,
	"POST":         PostHandler,
	"STAT":         StatHandler,
	"QUIT":         QuitHandler,
	"STARTTLS":     StarttlsHandler,
}

// StorageBackend is an alias for retrieving the storage interface associated
// with the server that accepted this connection
func (c *Conn) StorageBackend() Storage {
	return c.server.storage
}

// MessageFilter is an alias for retrieving the message filtering function associated
// with the server that accepted this connection
func (c *Conn) MessageFilter() FilterFunc {
	return c.server.filter
}

// AuthBackend is an alias for retrieving the authentication interface associated
// with the server that accepted this connection
func (c *Conn) AuthBackend() Auth {
	return c.server.auth
}

// CurrentArticle retrieves the article pointed to by the connections current group
// and article number, returning nil if there is no such existing article
func (c *Conn) CurrentArticle() *Article {
	if c.articleNumber != nil {
		s := c.StorageBackend()
		a, err := s.ArticleByGroup(*c.group, *c.articleNumber)
		if err != nil {
			return nil
		} else {
			return a
		}
	} else {
		return nil
	}
}

// ReadLine reads a CR-LF delimited line from the socket
func (c *Conn) ReadLine() (string, error) {
	reader := textproto.NewReader(c.br)
	return reader.ReadLine()
}

// ReadArticle reads a dot-encoded, CR-LF delimited MIME message from the socket
func (c *Conn) ReadArticle() (*Article, error) {
	reader := textproto.NewReader(c.br)
	header, err := reader.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	return &Article{
		header,
		reader.DotReader(),
	}, nil

}

// WriteLine formats a string and writes it to the socket with CR-LF line ending
func (c *Conn) WriteLine(text string, args ...interface{}) error {
	return textproto.NewWriter(c.bw).PrintfLine(text, args...)
}

func (c *Conn) WriteResponse(code int, args ...interface{}) error {
	return textproto.NewWriter(c.bw).PrintfLine(ResponseText(code, args...))
}

// WriteHeaders writes a dot-encoded listing of article headers to the socket with CR-LF delimiters
func (c *Conn) WriteHeaders(article Article) error {
	writer := textproto.NewWriter(c.bw).DotWriter()
	h := http.Header(article.MIMEHeader)
	if err := h.Write(writer); err != nil {
		return err
	}
	return writer.Close()
}

// WriteBody writes a dot-encoded CR-LF delimited message to the socket
func (c *Conn) WriteBody(article Article) error {
	writer := textproto.NewWriter(c.bw).DotWriter()
	if _, err := io.Copy(writer, article.Body); err != nil {
		return err
	}
	return writer.Close()
}

// WriteArticle writes a dot-encoded CR-LF delimited MIME message to the socket
func (c *Conn) WriteArticle(article Article) error {
	writer := textproto.NewWriter(c.bw).DotWriter()
	if err := http.Header(article.MIMEHeader).Write(writer); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, ""); err != nil {
		return err
	}
	if _, err := io.Copy(writer, article.Body); err != nil {
		return err
	}
	return writer.Close()
}
