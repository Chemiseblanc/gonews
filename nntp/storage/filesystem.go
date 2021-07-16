package storage

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Chemiseblanc/gonews/nntp"
)

type LegacyFileSystem struct{}

func (l *LegacyFileSystem) HasArticle(id string) (bool, error) {
	article, err := l.GetArticleByID(id)
	if err != nil {
		return false, err
	}

	if article != nil {
		return true, nil
	} else {
		return false, nil
	}
}

func (*LegacyFileSystem) GetGroup(group string) (*nntp.Group, error) {
	active, err := os.Open("/var/spool/news/active")
	if err != nil {
		return nil, err
	}

	found := false
	var description string
	var min uint
	var max uint
	var flag string

	scanner := bufio.NewScanner(active)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			return nil, nil
		}
		if fields[0] == group {
			found = true
			min, _ := strconv.ParseUint(fields[1], 10, 0)
			max, _ := strconv.ParseUint(fields[2], 10, 0)
			flag = fields[3]
			break
		}

	}
	return nil, nil
}

func (*LegacyFileSystem) PostArticle(article nntp.Article) error {
	return nil
}

func (*LegacyFileSystem) GetArticleByID(id string) (*nntp.Article, error) {
	history, err := os.Open("/var/spool/news/history")
	if err != nil {
		return nil, err
	}
	hash := fmt.Sprintf("%X", sha1.Sum([]byte(id)))

	scanner := bufio.NewScanner(history)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			return nil, nil
		}
		if fields[0] == hash {

		}
	}

}

func (*LegacyFileSystem) GetArticleByGroup(group nntp.Group, number uint) (*nntp.Article, error) {

}
