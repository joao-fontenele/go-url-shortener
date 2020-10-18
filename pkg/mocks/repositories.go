package mocks

import (
	"context"

	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

// FakeLinkRepo holds fake implementations for the LinkRepository interface
type FakeLinkRepo struct {
	ListFn     func(ctx context.Context, limit, skip int) ([]shortener.Link, error)
	ListCalled bool

	FindFn     func(ctx context.Context, slug string) (*shortener.Link, error)
	FindCalled bool

	DeleteFn     func(ctx context.Context, slug string) error
	DeleteCalled bool

	InsertFn     func(ctx context.Context, l *shortener.Link) (*shortener.Link, error)
	InsertCalled bool

	UpdateFn     func(ctx context.Context, l *shortener.Link) error
	UpdateCalled bool
}

// ensure FakeLinkRepo implements shortener.LinkRepository
var _ shortener.LinkRepository = &FakeLinkRepo{}

// Find is a mock for Find method in link repository
func (lr *FakeLinkRepo) Find(ctx context.Context, slug string) (*shortener.Link, error) {
	lr.FindCalled = true
	return lr.FindFn(ctx, slug)
}

// Delete is a mock for Delete method in link repository
func (lr *FakeLinkRepo) Delete(ctx context.Context, slug string) error {
	lr.DeleteCalled = true
	return lr.DeleteFn(ctx, slug)
}

// Insert is a mock for Insert method in link repository
func (lr *FakeLinkRepo) Insert(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
	lr.InsertCalled = true
	return lr.InsertFn(ctx, l)
}

// Update is a mock for Update method in link repository
func (lr *FakeLinkRepo) Update(ctx context.Context, l *shortener.Link) error {
	lr.UpdateCalled = true
	return lr.UpdateFn(ctx, l)
}

// List is a mock for List method in link repository
func (lr *FakeLinkRepo) List(ctx context.Context, limit, skip int) ([]shortener.Link, error) {
	lr.ListCalled = true
	return lr.ListFn(ctx, limit, skip)
}
