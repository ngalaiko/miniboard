package http

import (
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"miniboard.app/reader"
)

// Reader returns simplified HTML content.
type Reader struct {
	root *html.Node
	url  *url.URL
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
	n, err := html.Parse(raw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse html")
	}
	return &Reader{
		root: n,
		url:  url,
	}, nil
}

// Title returns the page title.
func (r *Reader) Title() (title string) {
	bfs(r.root, func(n *html.Node) bool {
		if n.DataAtom != atom.Head {
			return true
		}
		bfs(n, func(n *html.Node) bool {
			if n.DataAtom != atom.Title {
				return true
			}

			if n.FirstChild == nil {
				return false
			}

			title = n.FirstChild.Data
			return false
		})
		return false
	})
	return
}

// IconURL returns a link to the first page favicon.
func (r *Reader) IconURL() (iconURLs []string) {
	bfs(r.root, func(n *html.Node) bool {
		if n.DataAtom != atom.Head {
			return true
		}

		bfs(n, func(n *html.Node) bool {
			if n.DataAtom != atom.Link {
				return true
			}
			for _, attr := range n.Attr {
				extention := filepath.Ext(attr.Val)
				if extention != ".png" && extention != ".ico" {
					continue
				}
				u, err := url.Parse(attr.Val)
				if err != nil {
					continue
				}
				if u.Hostname() != "" {
					iconURLs = append(iconURLs, u.String())
					return false
				}

				iconURLs = append(iconURLs, r.url.String()+attr.Val)
				return false
			}
			return false
		})
		return true
	})
	return
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

// dfs
