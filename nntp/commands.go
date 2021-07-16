package nntp

import (
	"bufio"
	"crypto/tls"
	"net"
	"strconv"
)

// isMessageID is a helper function for checking if a given argument refers to an article number of message-id
func isMessageID(identifier string) bool {
	if len(identifier) > 0 && identifier[0] == '<' {
		return true
	} else {
		return false
	}
}

// retrievalHandler is a unified implementation of the shared behaviour of ARTICLE, HEAD, and BODY commands
func retrievalHandler(c *Conn, args []string, responseHandler func(*Conn, uint, *Article) error) error {
	if len(args) > 1 {
		if err := c.WriteLine(ResponseText(ResponseCommandSyntaxError)); err != nil {
			return err
		}
		return nil
	}

	if len(args) == 1 {
		s := c.StorageBackend()

		// First or second form of ARTICLE command
		ident := args[0]
		if isMessageID(ident) {
			// Article specified by Message-ID
			article, err := s.ArticleByID(ident)
			if err != nil {
				return err
			}

			if article != nil {
				return responseHandler(c, 0, article)
			} else {
				return c.WriteLine(ResponseText(ResponseArticleNotFound))
			}
		} else {
			// Article specified by number in group
			g := c.group
			if g != nil {
				article_number, err := strconv.ParseUint(ident, 10, 0)
				if err != nil {
					return err
				}

				article, err := s.ArticleByGroup(*g, uint(article_number))
				if err != nil {
					return err
				}

				if article != nil {
					*c.articleNumber = uint(article_number)
					return responseHandler(c, uint(article_number), article)
				} else {
					return c.WriteLine(ResponseText(ResponseArticleNotFound))
				}
			} else {
				return c.WriteLine(ResponseText(ResponseGroupNotSelected))
			}
		}
	} else {
		// Third form of article command
		g := c.group
		if g != nil {
			article := c.CurrentArticle()
			if article != nil {
				return responseHandler(c, *c.articleNumber, article)
			} else {
				return c.WriteLine(ResponseText(ResponseArticleNotSelected))
			}
		} else {
			return c.WriteLine(ResponseText(ResponseGroupNotSelected))
		}
	}
}

// Implements the ARTICLE command as described in section 6.2.1 of RFC3977
func ArticleHander(c *Conn, args []string) error {
	return retrievalHandler(c, args, func(c *Conn, number uint, a *Article) error {
		if err := c.WriteLine(ResponseText(ResponseArticleRetrievedHeadBody, number, a.MessageID())); err != nil {
			return err
		}
		return c.WriteArticle(*a)
	})
}

// Implements the BODY command as described in section 6.2.3 of RFC3977
func BodyHandler(c *Conn, args []string) error {
	return retrievalHandler(c, args, func(c *Conn, number uint, a *Article) error {
		if err := c.WriteLine(ResponseText(ResponseArticleRetrievedBody, number, a.MessageID())); err != nil {
			return err
		}
		return c.WriteBody(*a)
	})
}

// Implements the CAPABILITIES command as described in section 5.2 of RFC3977
func CapabilitiesHandler(c *Conn, args []string) error {
	caps := []string{
		"VERSION 2",
		"BODY",
		"CAPABILITIES",
		"DATE",
		"GROUP",
		"HDR",
		"HEAD",
		"HELP",
		"IHAVE",
		"LAST",
		"LIST",
		"LISTGROUP",
		"MODE",
		"NEWGROUPS",
		"NEWNEWS",
		"NEXT",
		"OVER",
		"POST",
		"STAT",
		"QUIT",
		"STARTTLS",
	}
	if err := c.WriteLine(ResponseText(ResponseCapabilitiesFollows)); err != nil {
		return err
	}
	for _, v := range caps {
		if err := c.WriteLine(v); err != nil {
			return err
		}
	}
	return c.WriteLine(".")
}

// Implements the DATE command as described in section 7.1 of RFC3977
func DateHandler(c *Conn, args []string) error {
	return nil
}

// Implements the GROUP command as described in section 6.1.1 of RFC3977
func GroupHandler(c *Conn, args []string) error {
	if len(args) != 1 {
		if err := c.WriteLine(ResponseText(ResponseCommandSyntaxError)); err != nil {
			return err
		}
		return nil
	}

	s := c.StorageBackend()
	g := s.Group(args[0])

	if g != nil {
		c.group = g
		if err := c.WriteLine(ResponseText(ResponseGroupSelected, g.Min, g.Max, g.Count, g.Name)); err != nil {
			return err
		}
	} else {
		if err := c.WriteLine(ResponseText(ResponseGroupNotFound)); err != nil {
			return err
		}
	}
	return nil
}

// Implements the HDR command as described in section 8.5 of RFC3977
func HdrHandler(c *Conn, args []string) error {
	return nil
}

// Implements the HEAD command as described in section 6.2.2 of RFC3977
func HeadHandler(c *Conn, args []string) error {
	return retrievalHandler(c, args, func(c *Conn, number uint, a *Article) error {
		if err := c.WriteLine(ResponseText(ResponseArticleRetrievedHead, number, a.MessageID())); err != nil {
			return err
		}
		return c.WriteHeaders(*a)
	})
}

// Implements the HELP command as described in section 7.2 of RFC3977
func HelpHandler(c *Conn, args []string) error {
	return nil
}

// Implements the IHAVE command as described in section 6.3.2 of RFC3977
func IhaveHandler(c *Conn, args []string) error {
	return nil
}

// Implements the LAST command as described in section 6.1.3 of RFC3977
func LastHandler(c *Conn, args []string) error {
	if len(args) > 0 {
		if err := c.WriteLine(ResponseText(ResponseCommandSyntaxError)); err != nil {
			return err
		}
		return nil
	}
	if g := c.group; g != nil {
		if *c.articleNumber >= g.Min {
			s := c.StorageBackend()
			for number := *c.articleNumber - 1; number >= g.Min; number-- {
				if a, err := s.ArticleByGroup(*g, number); err != nil {
					return err
				} else if a != nil {
					*c.articleNumber = number
					return c.WriteLine(ResponseText(ResponseArticleRetrieved, number, a.MessageID()))
				} else {
					continue
				}
			}
			return c.WriteLine(ResponseText(ResponseArticleNoPrevious))
		} else {
			return c.WriteLine(ResponseText(ResponseArticleNoPrevious))
		}
	} else {
		return c.WriteLine(ResponseText(ResponseGroupNotSelected))
	}
}

// Implements the LIST command as described in section 7.6.1 of RFC3977
func ListHandler(c *Conn, args []string) error {
	return nil
}

// Implements the LISTGROUP command as described in section 6.1.2 of RFC3977
func ListgroupHandler(c *Conn, args []string) error {
	return nil
}

// Implements the MODE command as described in section 5.3 of RFC3977
func ModeHandler(c *Conn, args []string) error {
	return nil
}

// Implements the NEWGROUPS command as described in section 7.3 of RFC3977
func NewgroupsHandler(c *Conn, args []string) error {
	return nil
}

// Implements the NEWNEWS command as described in section 7.4 of RFC3977
func NewnewsHandler(c *Conn, args []string) error {
	return nil
}

// Implements the NEXT command as described in section 6.1.4 of RFC3977
func NextHandler(c *Conn, args []string) error {
	if len(args) > 0 {
		if err := c.WriteLine(ResponseText(ResponseCommandSyntaxError)); err != nil {
			return err
		}
		return nil
	}
	if g := c.group; g != nil {
		if *c.articleNumber <= g.Max {
			s := c.StorageBackend()
			for number := *c.articleNumber + 1; number <= g.Max; number++ {
				if a, err := s.ArticleByGroup(*g, number); err != nil {
					return err
				} else if a != nil {
					*c.articleNumber = number
					return c.WriteLine(ResponseText(ResponseArticleRetrieved, number, a.MessageID()))
				} else {
					continue
				}
			}
			return c.WriteLine(ResponseText(ResponseArticleNoNext))
		} else {
			return c.WriteLine(ResponseText(ResponseArticleNoNext))
		}
	} else {
		return c.WriteLine(ResponseText(ResponseGroupNotSelected))
	}
}

// Implements the OVER command as described in section 8.3 of RFC3977
func OverHandler(c *Conn, args []string) error {
	return nil
}

// Implements the POST command as described in section 6.3.1 of RFC3977
func PostHandler(c *Conn, args []string) error {
	if len(args) > 0 {
		if err := c.WriteLine(ResponseText(ResponseCommandSyntaxError)); err != nil {
			return err
		}
		return nil
	}

	if c.AuthBackend().AnonymousPostingAllowed() {
		if err := c.WriteLine(ResponseText(ResponsePostArticle)); err != nil {
			return err
		}
		article, err := c.ReadArticle()
		if err != nil {
			c.WriteLine(ResponseText(ResponseArticleTransferFailed))
			return err
		}
		if article != nil {
			f := c.MessageFilter()
			if !f(*article) {
				s := c.StorageBackend()
				s.PostArticle(*article)
				return c.WriteLine(ResponseText(ResponseArticlePosted))
			} else {
				return c.WriteLine(ResponseText(ResponsePostingFailed))
			}
		} else {
			return c.WriteLine(ResponseText(ResponseArticleTransferFailed))
		}
	} else {
		return c.WriteLine(ResponseText(ResponsePostingNotAllowed))
	}
}

// Implements the STAT command as described in section 6.2.4 of RFC3977
func StatHandler(c *Conn, args []string) error {
	return retrievalHandler(c, args, func(c *Conn, number uint, a *Article) error {
		return c.WriteLine(ResponseText(ResponseArticleRetrieved, number, a.MessageID()))
	})
}

// Implements the QUIT command as described in section 6.2.4 of RFC3977
func QuitHandler(c *Conn, args []string) error {
	if err := c.WriteLine(ResponseText(ResponseConnectionClosing)); err != nil {
		return err
	}
	return c.Close()
}

// Implements the STARTTLS command as described in section 2.2 of RFC4642
func StarttlsHandler(c *Conn, args []string) error {
	if c.server.TLSConfig != nil {
		tlsConn := tls.Server(c.Conn, c.server.TLSConfig)
		c.br = bufio.NewReader(tlsConn)
		c.bw = bufio.NewWriter(tlsConn)
		c.Conn = net.Conn(tlsConn)
		c.isTLS = true
		return nil
	} else {
		return c.WriteLine(ResponseText(ResponseCommandNotSupported))
	}
}
