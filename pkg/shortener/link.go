package shortener

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Link holds the attributes related to shortened link urls
type Link struct {
	Slug      string    `json:"slug"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

// LinkDao represents a contract to access a single datastore
type LinkDao interface {
	Find(ctx context.Context, slug string) (*Link, error)
	Insert(ctx context.Context, l *Link) (*Link, error)
	Update(ctx context.Context, l *Link) error
	Delete(ctx context.Context, slug string) error
}

// Validate checks if a link is valid
func (l *Link) Validate() error {
	if l == nil {
		return fmt.Errorf("%w: Link should not be nil", ErrInvalidLink)
	}

	u, err := url.Parse(l.URL)

	if err != nil {
		return fmt.Errorf("%w: Parsing Link.URL generated error", ErrInvalidLink)
	}

	if u.Host == "" || u.Scheme == "" {
		return fmt.Errorf("%w: Link URL is malformed", ErrInvalidLink)
	}

	return err
}
