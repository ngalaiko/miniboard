package reader

// Reader returns article's info.
type Reader interface {
	Title() string
	IconURL() []string
}
