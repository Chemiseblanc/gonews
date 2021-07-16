package nntp

import (
	"io"
	"net/textproto"
)

type State struct {
	article uint
	group   string
	first   uint
	last    uint
}

type Article struct {
	id      string
	headers map[string]string
	body    string
}

type Storage interface {
	HasArticle(string) bool
	PostArticle(Article) error
	GetArticleID(string) (Article, error)
	GetArticleGroup(string, uint) (Article, error)
}

type Auth interface {
	AnonymousPostingAllowed() bool
}

type Filter interface{}
