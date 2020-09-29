package mocks

import (
	"context"

	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

// FakeLinkService holds fake implementations for the LinkService interface
type FakeLinkService struct {
	GetURLFn     func(ctx context.Context, slug string) (string, error)
	GetURLCalled bool

	CreateFn     func(ctx context.Context, URL string) (*shortener.Link, error)
	CreateCalled bool

	GetNewSlugFn     func(ctx context.Context, size int) (string, error)
	GetNewSlugCalled bool

	GenerateSlugFn     func(size int) string
	GenerateSlugCalled bool
}

// ensures FakeLinkService implements LinkService interface
var _ shortener.LinkService = &FakeLinkService{}

// GetURL returns an URL given a shortened url
func (ls *FakeLinkService) GetURL(ctx context.Context, slug string) (string, error) {
	ls.GetURLCalled = true
	return ls.GetURLFn(ctx, slug)
}

// Create creates a searchable URL for a given code
func (ls *FakeLinkService) Create(ctx context.Context, URL string) (*shortener.Link, error) {
	ls.CreateCalled = true
	return ls.CreateFn(ctx, URL)
}

// GetNewSlug returns a slug that still doesn't exist in db
func (ls *FakeLinkService) GetNewSlug(ctx context.Context, size int) (string, error) {
	ls.GetNewSlugCalled = true
	return ls.GetNewSlugFn(ctx, size)
}

// GenerateSlug returns a random slug
func (ls *FakeLinkService) GenerateSlug(size int) string {
	ls.GenerateSlugCalled = true
	return ls.GenerateSlugFn(size)
}
