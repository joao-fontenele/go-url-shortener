package shortener_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/joao-fontenele/go-url-shortener/pkg/mocks"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

func TestGenerateSlug(t *testing.T) {
	s := shortener.NewLinkService(&mocks.FakeLinkRepo{})
	slug1 := s.GenerateSlug(5)
	if len(slug1) != 5 {
		t.Errorf("Expected generate slug to have length 5, but got slug = %s", slug1)
	}

	slug2 := s.GenerateSlug(5)
	if len(slug2) != 5 {
		t.Errorf("Expected generate slug to have length 5, but got slug = %s", slug2)
	}

	if slug1 == slug2 {
		t.Errorf("Generated slugs should have been random but got %s = %s", slug1, slug2)
	}
}

func TestGetNewSlug(t *testing.T) {
	t.Run("EventualSuccess", func(t *testing.T) {
		attempts := 0
		s := shortener.NewLinkService(&mocks.FakeLinkRepo{
			FindFn: func(ctx context.Context, slug string) (*shortener.Link, error) {
				attempts++
				if attempts <= 4 {
					return &shortener.Link{}, nil
				}
				return nil, shortener.ErrLinkNotFound
			},
		})

		slug, err := s.GetNewSlug(context.Background(), 5)

		if err != nil {
			t.Fatalf("Unexpected error in GetNewSlug: %v", err)
		}

		if len(slug) != 5 {
			t.Errorf("Slug should have length 5, but got slug %s", slug)
		}

		if attempts != 5 {
			t.Errorf("GetNewSlug should have errored 4 times before succeeding, attempts = %d", attempts)
		}
	})

	t.Run("Failure", func(t *testing.T) {
		fakeRepo := mocks.FakeLinkRepo{
			FindFn: func(ctx context.Context, slug string) (*shortener.Link, error) {
				return &shortener.Link{}, nil
			},
		}
		s := shortener.NewLinkService(&fakeRepo)
		_, err := s.GetNewSlug(context.Background(), 5)

		if !errors.Is(err, shortener.ErrLinkExists) {
			t.Fatalf("Expected ErrLinkExists, but got: %v", err)
		}

		if !fakeRepo.FindCalled {
			t.Errorf("Expected Find to have been called, but it wasn't called")
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("Failure", func(t *testing.T) {
		fakeRepo := mocks.FakeLinkRepo{
			FindFn: func(ctx context.Context, slug string) (*shortener.Link, error) {
				return &shortener.Link{}, nil
			},
		}

		s := shortener.NewLinkService(&fakeRepo)

		_, err := s.Create(context.Background(), "https://www.google.com")
		if err == nil {
			t.Errorf("Expected error to not be nil, but got: %v", err)
		}

		if !fakeRepo.FindCalled {
			t.Errorf("Expected Find to have been called, but it wasn't called")
		}
	})

	t.Run("Success", func(t *testing.T) {
		fakeRepo := &mocks.FakeLinkRepo{
			FindFn: func(ctx context.Context, slug string) (*shortener.Link, error) {
				return nil, shortener.ErrLinkNotFound
			},
			InsertFn: func(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
				l.CreatedAt = time.Now()
				return l, nil
			},
		}
		s := shortener.NewLinkService(fakeRepo)

		link, err := s.Create(context.Background(), "https://www.google.com")

		if err != nil {
			t.Fatalf("Unexpected error while creating Link: %v", err)
		}

		if !fakeRepo.FindCalled {
			t.Errorf("Expected Find to have been called, but it wasn't called")
		}

		if !fakeRepo.InsertCalled {
			t.Errorf("Expected Insert to have been called, but it wasn't called")
		}

		if link.URL != "https://www.google.com" {
			t.Errorf("Expected Link.URL to be 'https://www.google.com', but got: %s", link.URL)
		}

		if len(link.Slug) != 5 {
			t.Errorf("Expected Link.Slug to have len = 5, but got: %s", link.Slug)
		}
	})
}

func TestGetURL(t *testing.T) {
	t.Run("LinkFound", func(t *testing.T) {
		fakeRepo := mocks.FakeLinkRepo{
			FindFn: func(ctx context.Context, slug string) (*shortener.Link, error) {
				return &shortener.Link{URL: "https:/www.google.com", Slug: "dummy"}, nil
			},
		}

		s := shortener.NewLinkService(&fakeRepo)
		URL, err := s.GetURL(context.Background(), "dummy")

		if err != nil {
			t.Fatalf("Unexpected error from GetURL: %v", err)
		}

		if URL != "https:/www.google.com" {
			t.Errorf("Expected URL to be 'https://www.google.com', but got: %s", URL)
		}
	})

	t.Run("LinkNotFound", func(t *testing.T) {
		fakeRepo := mocks.FakeLinkRepo{
			FindFn: func(ctx context.Context, slug string) (*shortener.Link, error) {
				return nil, shortener.ErrLinkNotFound
			},
		}

		s := shortener.NewLinkService(&fakeRepo)
		URL, err := s.GetURL(context.Background(), "dummy")

		if !errors.Is(err, shortener.ErrLinkNotFound) {
			t.Fatalf("Unexpected error from GetURL: %v", err)
		}

		if URL != "" {
			t.Errorf("Expected URL to be '', but got: %s", URL)
		}
	})
}
