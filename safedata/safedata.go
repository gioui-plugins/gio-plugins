package safedata

import "errors"

var (
	// DefaultAppName is the unique name to identify the app,
	// you should change it, or provide one name using Config.
	//
	// The DefaultAppName is used when Config.App is empty.
	DefaultAppName = "safedata"
)

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

	// ErrMetadataMaxLength is returned when the given identifier
	// hits or exceeds the maximum allowed length.
	ErrMetadataMaxLength = errors.New("metadata exceeds the maximum length")

	// ErrMalformedMetadata is returned when the given identifier
	// contains invalid character. Usually that happens when identifier
	// contains control characters (null, tabs...).
	ErrMalformedMetadata = errors.New("metadata contains invalid characters")
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
	// That field is not encrypted.
	Identifier string

	// Description is a brief description explaining
	// this data. It may be displayed to the end-user,
	// describing the purpose of that credential. It
	// might be a good practices to include the name
	// of your software into the description.
	//
	// That field is ignored on Android and WASM.
	// That field is not encrypted.
	Description string

	// Data is the arbitrary secret data, it can be one
	// password, token or certificate.
	//
	// That field is encrypted on Android, iOS,
	// macOS and Windows.
	Data []byte
}

type SafeData struct {
	driver
}

func NewSafeData(config Config) *SafeData {
	if config.App == "" {
		config.App = DefaultAppName
	}

	h := SafeData{}
	attachDriver(&h, config)
	return &h
}

type Looper func(identifier string) (next bool)

type safeData interface {
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
func (s *SafeData) Set(secret Secret) error {
	return s.setSecret(secret)
}

// Get gets the data from credentials manager.
//
// The identifier must match against the identifier previously
// used for Set.
func (s *SafeData) Get(identifier string) (Secret, error) {
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
func (s *SafeData) View(identifier string, out *Secret) error {
	return s.getSecret(identifier, out)
}

// List gets a list of credentials that belongs to the
// current app.
//
// On WebAssembly that will always return a maximum of one
// credential, if any.
func (s *SafeData) List(looper Looper) error {
	return s.listSecret(looper)
}

// Remove deletes the data and identifier from OS
// credentials manager.
//
// The identifier is used for searching.
func (s *SafeData) Remove(identifier string) error {
	return s.removeSecret(identifier)
}

func (s *SafeData) Configure(config Config) {
	if config.App == "" {
		config.App = DefaultAppName
	}

	configureDriver(&s.driver, config)
}
