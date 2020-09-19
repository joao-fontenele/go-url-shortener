package shortener

import (
	"context"
	"time"
)

// Link holds the attributes related to shortened link urls
type Link struct {
	Slug      string
	URL       string
	CreatedAt time.Time
}

// LinkDao represents a contract to access a single datastore
type LinkDao interface {
	Find(ctx context.Context, slug string) (*Link, error)
	Insert(ctx context.Context, l *Link) (*Link, error)
	Update(ctx context.Context, l *Link) error
	Delete(ctx context.Context, slug string) error
}
