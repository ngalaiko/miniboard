package reader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-shiori/go-readability"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"miniboard.app/images"
	"miniboard.app/storage/resource"
)

// GetClient is used to fetch article's data from the Internet.
type GetClient interface {
	Get(string) (*http.Response, error)
}

// Reader returns simplified HTML content.
type Reader struct {
	article *readability.Article
	url     *url.URL
}

// New returns new reader from a url.
func New(ctx context.Context, client GetClient, articleName *resource.Name, images *images.Service, url *url.URL) (*Reader, error) {
	resp, err := client.Get(url.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch url")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return NewFromReader(ctx, client, articleName, images, resp.Body, url)
}

// NewFromReader returns new reader from io.Reader.
// URL is needed to form complete links to images.
func NewFromReader(ctx context.Context, client GetClient, articleName *resource.Name, images *images.Service, raw io.Reader, url *url.URL) (*Reader, error) {
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
			inlineImage(ctx, client, articleName, images, n)
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

func inlineImage(ctx context.Context, client GetClient, articleName *resource.Name, images *images.Service, n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key != "src" {
			continue
		}

		name, err := images.SaveURL(ctx, articleName, attr.Val)
		if err != nil {
			return
		}

		n.Attr = []html.Attribute{
			{
				Namespace: attr.Namespace,
				Key:       "src",
				Val:       fmt.Sprintf("/%s", name),
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
