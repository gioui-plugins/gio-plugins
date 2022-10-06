package installview

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Installer is the WebView installer.
type Installer struct {
	Data []byte
}

// Install will install the Installer.
//
// That will vary between OSes.
func (i *Installer) Install() error {
	return i.install()
}

// Download will download the Installer.
//
// The download link will vary between OSes.
func Download(ctx context.Context, client *http.Client) (*Installer, error) {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DownloadURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Installer{Data: data}, nil
}
