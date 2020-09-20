package mocks

import (
	"context"

	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

// FakeLinkDao holds fake implementations for the LinkDao interface
type FakeLinkDao struct {
	FindFn     func(ctx context.Context, slug string) (*shortener.Link, error)
	FindCalled bool

	DeleteFn     func(ctx context.Context, slug string) error
	DeleteCalled bool

	InsertFn     func(ctx context.Context, l *shortener.Link) (*shortener.Link, error)
	InsertCalled bool

	UpdateFn     func(ctx context.Context, l *shortener.Link) error
	UpdateCalled bool
}

// ensure FakeLinkDao implements shortener.LinkDao
var _ shortener.LinkDao = &FakeLinkDao{}

// Find is a mock for Find method in link repository
func (lr *FakeLinkDao) Find(ctx context.Context, slug string) (*shortener.Link, error) {
	lr.FindCalled = true
	return lr.FindFn(ctx, slug)
}

// Delete is a mock for Delete method in link repository
func (lr *FakeLinkDao) Delete(ctx context.Context, slug string) error {
	lr.DeleteCalled = true
	return lr.DeleteFn(ctx, slug)
}

// Insert is a mock for Insert method in link repository
func (lr *FakeLinkDao) Insert(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
	lr.InsertCalled = true
	return lr.InsertFn(ctx, l)
}

// Update is a mock for Update method in link repository
func (lr *FakeLinkDao) Update(ctx context.Context, l *shortener.Link) error {
	lr.UpdateCalled = true
	return lr.UpdateFn(ctx, l)
}
