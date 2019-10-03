package resource

import (
	"bytes"
	"strings"
)

// Name is a resource name.
type Name struct {
	// subresource ids ordered by ownership.
	parts []*idPart
}

// id of a single resource
type idPart struct {
	ID   string
	Type string
}

// AddChild creates a new name with child.
func (n *Name) AddChild(child *Name) *Name {
	n.parts = append(n.parts, child.parts...)
	return n
}

// NewName creates a new id for the _typ_.
func NewName(typ string, id string) *Name {
	part := &idPart{
		ID:   id,
		Type: typ,
	}

	return &Name{
		parts: []*idPart{part},
	}
}

// ParseName parses a name from `type/id` string.
func ParseName(str string) *Name {
	parts := strings.Split(str, "/")

	name := &Name{
		parts: make([]*idPart, 0, len(parts)/2),
	}
	for i := 0; i < len(parts)-1; i += 2 {
		name.parts = append(name.parts, &idPart{
			ID:   parts[i+1],
			Type: parts[i],
		})
	}
	return name
}

// Path returns all resources in the path.
func (n *Name) Path() []*Name {
	names := make([]*Name, 0, len(n.parts))
	for _, p := range n.parts {
		names = append(names, &Name{
			parts: []*idPart{p},
		})
	}
	return names
}

// Split returns name splitted into id and everything else.
func (n *Name) Split() (string, string) {
	first := &strings.Builder{}
	for i, p := range n.parts {
		if i == len(n.parts)-1 {
			first.WriteString(p.Type)
			return first.String(), p.ID
		}
		first.WriteString(p.Type)
		first.WriteRune('/')
		first.WriteString(p.ID)
		first.WriteRune('/')
	}
	return "", ""
}

// Type returns resource's type.
func (n *Name) Type() string {
	return n.parts[len(n.parts)-1].Type
}

// Parent returns resource's parent.
func (n *Name) Parent() *Name {
	partsCopy := make([]*idPart, len(n.parts))
	copy(partsCopy, n.parts)
	return &Name{
		parts: n.parts[:len(partsCopy)-1],
	}
}

// ID returns resource's id.
func (n *Name) ID() string {
	return n.parts[len(n.parts)-1].ID
}

// Child returns a new child _resource_.
func (n *Name) Child(typ string, id string) *Name {
	partsCopy := make([]*idPart, len(n.parts))
	copy(partsCopy, n.parts)
	return &Name{
		parts: append(partsCopy, &idPart{
			ID:   id,
			Type: typ,
		}),
	}
}

const delimmer = byte('/')

// Bytes returns full name as a bytes slice.
func (n *Name) Bytes() []byte {
	buffer := &bytes.Buffer{}
	for i, p := range n.parts {
		buffer.WriteString(p.Type)
		buffer.WriteByte(delimmer)
		buffer.WriteString(p.ID)

		if i != len(n.parts)-1 {
			buffer.WriteByte(delimmer)
		}
	}

	return buffer.Bytes()
}

// String returns full name as a string.
func (n *Name) String() string {
	return string(n.Bytes())
}
