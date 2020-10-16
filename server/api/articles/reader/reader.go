package reader

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	readability "github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Reader wraps readability parser.
type Reader struct {
	parser readability.Parser
}

// New returns a new readability reader.
func New() *Reader {
	return &Reader{
		parser: readability.NewParser(),
	}
}

// Parse returns ewadable contant of a page.
func (r *Reader) Parse(raw io.Reader, url string) ([]byte, string, error) {
	article, err := r.parser.Parse(raw, url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse document: %w", err)
	}

	wg := &sync.WaitGroup{}
	bfs(article.Node, func(n *html.Node) bool {
		if n.DataAtom != atom.A {
			return true
		}

		// open links on a new page
		n.Attr = append(n.Attr, html.Attribute{
			Key: "target",
			Val: "_blank",
		})

		return true
	})

	wg.Wait()

	b := &bytes.Buffer{}
	if err := html.Render(b, article.Node); err != nil {
		return nil, "", err
	}

	return b.Bytes(), article.Title, nil
}

func bfs(node *html.Node, forEach func(*html.Node) bool) {
	if node == nil {
		return
	}
	if !forEach(node) {
		return
	}
	n := node.FirstChild
	for n != nil {
		bfs(n, forEach)
		n = n.NextSibling
	}
}
