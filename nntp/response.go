package nntp

import "fmt"

var Quiet bool = false

const (
	ResponseCapabilitiesFollows      = 101
	ResponseServerReadyPosting       = 200
	ResponseServerReadyNoPosting     = 201
	ResponseConnectionClosing        = 205
	ResponseGroupSelected            = 211
	ResponseGroupListFollows         = 215
	ResponseGroupNotFound            = 511
	ResponseGroupNotSelected         = 512
	ResponseArticleRetrievedHeadBody = 220
	ResponseArticleRetrievedHead     = 221
	ResponseArticleRetrievedBody     = 222
	ResponseArticleRetrieved         = 223
	ResponseArticleTransferred       = 235
	ResponseArticlePosted            = 240
	ResponseTransferArticle          = 335
	ResponsePostArticle              = 340
	ResponseArticleNotSelected       = 420
	ResponseArticleNoNext            = 421
	ResponseArticleNoPrevious        = 422
	ResponseArticleNotInGroup        = 423
	ResponseArticleNotFound          = 430
	ResponseArticleNotWanted         = 435
	ResponseArticleTransferFailed    = 436
	ResponseArticleRejected          = 437
	ResponsePostingNotAllowed        = 440
	ResponsePostingFailed            = 441
	ResponseCommandNotRecognized     = 500
	ResponseCommandSyntaxError       = 501
	ResponseCommandNotSupported      = 503
)

var responseText = map[int]string{
	ResponseCapabilitiesFollows:      "%d capability list follows (multi-line)",
	ResponseServerReadyPosting:       "%d server ready - posting allowed",
	ResponseServerReadyNoPosting:     "%d server ready - no posting allowed",
	ResponseConnectionClosing:        "%d closing connection - goodbye!",
	ResponseGroupSelected:            "%d %d %d %d %s group selected",
	ResponseGroupListFollows:         "%d list of newsgroups follows",
	ResponseGroupNotFound:            "%d no such news group",
	ResponseGroupNotSelected:         "%d no newsgroup has been selected",
	ResponseArticleRetrievedHeadBody: "%d %d %s article retrieved - head and body follow",
	ResponseArticleRetrievedHead:     "%d %d %s article retrieved - head follows",
	ResponseArticleRetrievedBody:     "%d %d %s article retrieved - body follows",
	ResponseArticleRetrieved:         "%d %d %s article retrieved - request text seperately",
	ResponseArticleTransferred:       "%d article transferred ok",
	ResponseArticlePosted:            "%d article posted ok",
	ResponseTransferArticle:          "%d send article to be transferred. End with <CR-LF>.<CR-LF>",
	ResponsePostArticle:              "%d send article to be posted. End with <CR-LF>.<CR-LF>",
	ResponseArticleNotSelected:       "%d no current article has been selected",
	ResponseArticleNoNext:            "%d no next article in this group",
	ResponseArticleNoPrevious:        "%d no previous article in this group",
	ResponseArticleNotInGroup:        "%d no such article number in this group",
	ResponseArticleNotFound:          "%d no such article found",
	ResponseArticleNotWanted:         "%d article not wanted - do not send it",
	ResponseArticleTransferFailed:    "%d transfer failed - try again later",
	ResponseArticleRejected:          "%d article rejected - do not try again",
	ResponsePostingNotAllowed:        "%d posting not allowed",
	ResponsePostingFailed:            "%d posting failed",
	ResponseCommandNotRecognized:     "%d command not recognized",
	ResponseCommandSyntaxError:       "%d command syntax error",
	ResponseCommandNotSupported:      "%d command not supported",
}

const (
	quietStatusCode       = "%d"
	quietGroupSelected    = "%d %d %d %d %s"
	quietArticleRetrieved = "%d %d %s"
)

var responseTextQuiet = map[int]string{
	ResponseCapabilitiesFollows:      quietStatusCode,
	ResponseServerReadyPosting:       quietStatusCode,
	ResponseServerReadyNoPosting:     quietStatusCode,
	ResponseConnectionClosing:        quietStatusCode,
	ResponseGroupSelected:            quietGroupSelected,
	ResponseGroupListFollows:         quietStatusCode,
	ResponseGroupNotFound:            quietStatusCode,
	ResponseGroupNotSelected:         quietStatusCode,
	ResponseArticleRetrievedHeadBody: quietArticleRetrieved,
	ResponseArticleRetrievedHead:     quietArticleRetrieved,
	ResponseArticleRetrievedBody:     quietArticleRetrieved,
	ResponseArticleRetrieved:         quietArticleRetrieved,
	ResponseArticleTransferred:       quietStatusCode,
	ResponseArticlePosted:            quietStatusCode,
	ResponseTransferArticle:          quietStatusCode,
	ResponsePostArticle:              quietStatusCode,
	ResponseArticleNotSelected:       quietStatusCode,
	ResponseArticleNoNext:            quietStatusCode,
	ResponseArticleNoPrevious:        quietStatusCode,
	ResponseArticleNotInGroup:        quietStatusCode,
	ResponseArticleNotFound:          quietStatusCode,
	ResponseArticleNotWanted:         quietStatusCode,
	ResponseArticleTransferFailed:    quietStatusCode,
	ResponseArticleRejected:          quietStatusCode,
	ResponsePostingNotAllowed:        quietStatusCode,
	ResponsePostingFailed:            quietStatusCode,
	ResponseCommandNotRecognized:     quietStatusCode,
	ResponseCommandSyntaxError:       quietStatusCode,
	ResponseCommandNotSupported:      quietStatusCode,
}

func ResponseText(code int, param ...interface{}) string {
	var format string
	if Quiet {
		format = responseTextQuiet[code]
	} else {
		format = responseText[code]
	}
	return fmt.Sprintf(format, append([]interface{}{code}, param)...)
}
