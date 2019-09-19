package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/go-shiori/go-readability"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"miniboard.app/reader"
)

// Reader returns simplified HTML content.
type Reader struct {
	article *readability.Article
	url     *url.URL
}

// New erturns new reader from a url.
func New(url *url.URL) (reader.Reader, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch url")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return NewFromReader(resp.Body, url)
}

// NewFromReader returns new reader from io.Reader.
// URL is needed to form complete links to images.
func NewFromReader(raw io.Reader, url *url.URL) (*Reader, error) {
	article, err := readability.FromReader(raw, url.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse document")
	}
	return &Reader{
		article: &article,
		url:     url,
	}, nil
}

// Title returns the page title.
func (r *Reader) Title() (title string) {
	return r.article.Title
}

// Content returns page content.
func (r *Reader) Content() []byte {
	bfs(r.article.Node, func(node *html.Node) bool {
		if node.DataAtom == atom.Img {
			node.Attr = append(node.Attr, html.Attribute{
				Key: "style",
				Val: "width: 100%; height: auto",
			})
		}
		if node.DataAtom == atom.Pre {
			node.Attr = append(node.Attr, html.Attribute{
				Key: "style",
				Val: "overflow: auto",
			})
		}
		return true
	})
	buf := &bytes.Buffer{}
	_ = html.Render(buf, r.article.Node)
	return buf.Bytes()
}

// IconURL returns a link to the first page favicon.
func (r *Reader) IconURL() string {
	return r.article.Favicon
}

// executes forEach function on every node, including the first one in BFS order.
// if forEach returns true, search continues.
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
