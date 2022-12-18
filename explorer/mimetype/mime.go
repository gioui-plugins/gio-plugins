package mimetype

import "io"

// MimeType represents a mime type.
type MimeType struct {
	// Extension is the file extension associated with the mime type (e.g. ".jpg").
	Extension string
	// Type is the type part of the mime type (e.g. "image").
	Type string
	// Subtype is the subtype part of the mime type (e.g. "jpeg").
	Subtype string
}

// String returns the mime type as a string.
func (m MimeType) String() string {
	return m.Type + "/" + m.Subtype
}

// WriteTo writes the mime type to the given writer.
func (m MimeType) WriteTo(v io.Writer) {
	v.Write([]byte(m.Type))
	v.Write([]byte{'/'})
	v.Write([]byte(m.Subtype))
}
