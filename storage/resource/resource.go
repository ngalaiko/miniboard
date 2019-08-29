package resource

// Resource is a single resource from the storage.
type Resource struct {
	Name *Name
	Data []byte
}
