package nntp

const (
	server_ready_posting        = "200 server ready - posting allowed"
	server_ready_no_posting     = "201 server ready - no posting allowed"
	slave_noted                 = "202 slave status noted"
	connection_closing          = "205 closing connection - goodbye!"
	group_selected              = "211 %d %d %d %s group selected"
	group_list_follows          = "215 list of newsgroups follows"
	group_not_found             = "411 no such news group"
	group_not_selected          = "412 no newsgroup has been selected"
	article_retrieved_head_body = "220 %d %s article retrieved - head and body follow"
	article_retrieved_head      = "221 %d %s article retrieved - head follows"
	article_retrieved_body      = "222 %d %s article retrieved - body follows"
	article_retrieved           = "223 %d %s article retrieved - request text seperately"
	article_transferred         = "235 article transferred ok"
	article_posted              = "240 article posted ok"
	transfer_article            = "335 send article to be transferred. End with <CR-LF>.<CR-LF>"
	send_article                = "340 send article to be posted. End with <CR-LF>.<CR-LF>"
	article_not_selected        = "420 no current article has been selected"
	article_no_next             = "421 no next article in this group"
	article_no_previous         = "422 no previous article in this group"
	article_not_in_group        = "423 no such article number in this group"
	article_not_found           = "430 no such article found"
	article_not_wanted          = "435 article not wanted - do not send it"
	article_transfer_failed     = "436 transfer failed - try again later"
	article_rejected            = "437 article rejected - do not try again"
	posting_not_allowed         = "440 posting not allowed"
	command_not_recognized      = "500 command not recognized"
	command_syntax_error        = "501 command syntax error"
	command_not_supported       = "503 command not supportedd"
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
