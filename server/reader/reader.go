package reader

// Reader returns article's info.
type Reader interface {
	Title() string
	SiteName() string
	IconURL() string
	Content() []byte
}
