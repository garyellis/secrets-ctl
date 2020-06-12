package vault

import (
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client is the vault client
type Client struct {
	*Transit
	*KV
	*vaultapi.Client
}

// NewClient returns a new vault client
func NewClient() (*Client, error) {
	// Default Config returns the default config plus environment variable overrides
	config := vaultapi.DefaultConfig()
	if config.Error != nil {
		return nil, config.Error
	}

	config.HttpClient.Transport.(*http.Transport).TLSHandshakeTimeout = 5 * time.Second
	newclient, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		&Transit{
			newclient,
		},
		&KV{
			newclient,
		},
		newclient,
	}, nil
}
