package safehouse

import "errors"

var (
	// ErrNotFound is returned when there's no credentials
	// for the given key.
	ErrNotFound = errors.New("not found")

	// ErrUserRefused is returned when the user refuse to
	// save or retrieve credentials.
	ErrUserRefused = errors.New("user refused")

	// ErrUnsupported is returned when the current OS doesn't
	// supports the current requested feature.
	ErrUnsupported = errors.New("os not supported")

	// ErrIdentifierMaxLength is returned when the given identifier
	// hits or exceeds the maximum allowed length.
	ErrIdentifierMaxLength = errors.New("identifier exceeds the maximum length")
)

// Secret represents a data to be stored, and
// its searchable using Identifier.
type Secret struct {
	// Identifier is a unique identifier of the secret,
	// it is used to search the secret.
	//
	// The identifier is plain-text and usually is
	// the username or email of the user, or other
	// public identifiable content.
	//
	// If the Identifier exceeds 256 bytes, it will
	// be limited to 256 bytes.
	Identifier string

	// Description is a brief description explaining
	// this data. It may be displayed to the end-user,
	// describing the purpose of that credential. It
	// might be a good practices to include the name
	// of your software into the description.
	Description string

	// Data is the arbitrary secret data, it can be one
	// password, token or certificate.
	Data []byte
}

type SafeHouse struct {
	driver
}

func NewSafeHouse(config Config) *SafeHouse {
	h := SafeHouse{}
	attachDriver(&h, config)
	return &h
}

type Looper func(identifier string) (next bool)

type safeHouse interface {
	setSecret(secret Secret) error
	listSecret(looper Looper) error
	getSecret(identifier string, secret *Secret) error
	removeSecret(identifier string) error
}

// Set uploads (or updates) the given secret to the OS
// credentials manager.
//
// The identifier is used for searching, and must be unique,
// otherwise will replace the previous value.
func (s *SafeHouse) Set(secret Secret) error {
	return s.setSecret(secret)
}

// Get gets the data from credentials manager.
//
// The identifier must match against the identifier previously
// used for Set.
func (s *SafeHouse) Get(identifier string) (Secret, error) {
	var secret Secret
	err := s.getSecret(identifier, &secret)
	return secret, err
}

// View gets the data from credentials manager,
// and sets the content into the provided out. It may
// re-use the same Data slice.
//
// The "out" argument MUST NOT be nil. If you want to
// allocate new Secret for each call, see Get function.
//
// The identifier must match against the identifier previously
// used for Set.
func (s *SafeHouse) View(identifier string, out *Secret) error {
	return s.getSecret(identifier, out)
}

// List gets all credentials associated with the current App.
//
// On WebAssembly that will always return a maximum of one
// credential, if any.
func (s *SafeHouse) List(looper Looper) error {
	return s.listSecret(looper)
}

// List gets a list of credentials that belongs to the
// current app.

// Remove deletes the data and identifier from OS
// credentials manager.
//
// The identifier is used for searching.
func (s *SafeHouse) Remove(identifier string) error {
	return s.removeSecret(identifier)
}
