package nntp

import (
	"io"
	"net/textproto"
)

// Group is a structure describing a newsgroup the server participates in
type Group struct {
	Name        string
	Description string
	Min         uint
	Max         uint
	Count       uint
	Flag        string
}

// Article is a structure describing a news article.
type Article struct {
	textproto.MIMEHeader

	// Body is a reader that points to the source of the body text. This is done so large messages don't
	// have to be read into memory before operating on articles so unless the body text is required for
	// something like a filtering function the text can be transferred directly from the connection to the
	// storage decreasing memory pressure of the server under load
	Body io.Reader
}

// MessageID is a convenience function for retrieving the contents of the MessageID header field
func (a *Article) MessageID() string {
	return a.Get("MessageID")
}

// Storage is an interface for operations against the articles served by the newsserver
type Storage interface {
	HasArticle(string) bool
	Group(string) *Group
	PostArticle(Article) error
	ArticleByID(string) (*Article, error)
	ArticleByGroup(Group, uint) (*Article, error)
}

// Auth is an interface for validating whether or not to permit actions taken by an active connection
type Auth interface {
	AnonymousPostingAllowed() bool
}

// FilterFunc is a type of function for determining if the newsserver should accept a posted or transferred article
// it returns true if the given article should be rejected
type FilterFunc func(Article) bool
