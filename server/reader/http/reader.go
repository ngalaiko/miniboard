package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

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

	wg := &sync.WaitGroup{}
	bfs(article.Node, func(n *html.Node) bool {
		if n.DataAtom != atom.Img {
			return true
		}

		wg.Add(1)
		go func(n *html.Node) {
			inlineImage(n)
			wg.Done()
		}(n)
		return true
	})

	wg.Wait()

	return &Reader{
		article: &article,
		url:     url,
	}, nil
}

func inlineImage(n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key != "src" {
			continue
		}

		// TODO: mock http client, not reader
		resp, err := http.Get(attr.Val)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		n.Attr = []html.Attribute{
			{
				Namespace: attr.Namespace,
				Key:       "src",
				Val: fmt.Sprintf(
					"data:%s;base64,%s",
					resp.Header.Get("Content-Type"),
					base64.StdEncoding.EncodeToString(data),
				),
			},
		}
		break
	}
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

// Title returns the page title.
func (r *Reader) Title() (title string) {
	return r.article.Title
}

// SiteName returns name of source website.
func (r *Reader) SiteName() string {
	return r.article.SiteName
}

// Content returns page content.
func (r *Reader) Content() []byte {
	b := &bytes.Buffer{}
	html.Render(b, r.article.Node)
	return b.Bytes()
}

// IconURL returns a link to the first page favicon.
func (r *Reader) IconURL() string {
	return r.article.Favicon
}
