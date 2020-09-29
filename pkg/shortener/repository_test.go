package shortener_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/joao-fontenele/go-url-shortener/pkg/mocks"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

type findFn func(ctx context.Context, slug string) (*shortener.Link, error)

func TestFind(t *testing.T) {
	sampleLink := &shortener.Link{
		Slug:      "aaaaa",
		URL:       "htttps://wwww.google.com",
		CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
	}

	hitFind := func(ctx context.Context, slug string) (*shortener.Link, error) {
		return sampleLink, nil
	}

	missFind := func(ctx context.Context, slug string) (*shortener.Link, error) {
		return nil, shortener.ErrLinkNotFound
	}

	t.Run("CacheHit", func(t *testing.T) {
		db := &mocks.FakeLinkDao{FindFn: hitFind}
		cache := &mocks.FakeLinkDao{FindFn: hitFind}

		r := shortener.NewLinkRepository(db, cache)

		link, err := r.Find(context.Background(), "dontcare")

		if err != nil {
			t.Errorf("Unexpected error calling repository Find: %v", err)
		}

		if db.FindCalled {
			t.Error("Expected db find to not have been called")
		}

		if !cache.FindCalled {
			t.Error("Expected cache find to have been called")
		}

		if diff := cmp.Diff(sampleLink, link); diff != "" {
			t.Errorf("Found link different from expected (-want +got):\n%s", diff)
		}
	})

	t.Run("CacheMissDbMiss", func(t *testing.T) {
		db := &mocks.FakeLinkDao{FindFn: missFind}
		cache := &mocks.FakeLinkDao{FindFn: missFind}

		r := shortener.NewLinkRepository(db, cache)

		link, err := r.Find(context.Background(), "dontcare")

		if link != nil {
			t.Errorf("Expected link to be nil, but got: %v", link)
		}

		if !errors.Is(err, shortener.ErrLinkNotFound) {
			t.Errorf("Expected error to be ErrLinkNotFound, but got: %v", err)
		}

		if !db.FindCalled {
			t.Error("Expected db find to have been called")
		}

		if !cache.FindCalled {
			t.Error("Expected cache find to have been called")
		}
	})

	t.Run("CacheMissDbHit", func(t *testing.T) {
		db := &mocks.FakeLinkDao{FindFn: hitFind}
		cache := &mocks.FakeLinkDao{FindFn: missFind}

		r := shortener.NewLinkRepository(db, cache)

		link, err := r.Find(context.Background(), "dontcare")

		if err != nil {
			t.Errorf("Unexpected error calling repository Find: %v", err)
		}

		if !db.FindCalled {
			t.Error("Expected db find to not have been called")
		}

		if !cache.FindCalled {
			t.Error("Expected cache find to have been called")
		}

		if diff := cmp.Diff(sampleLink, link); diff != "" {
			t.Errorf("Found link different from expected (-want +got):\n%s", diff)
		}
	})
}

func TestInsert(t *testing.T) {
	sampleLink := &shortener.Link{
		URL:       "https://www.google.com/?search=Google",
		Slug:      "aaaaa",
		CreatedAt: time.Now(),
	}
	okInsert := func(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
		return sampleLink, nil
	}

	failInsert := func(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
		return nil, errors.New("UnexpectedErr")
	}

	t.Run("InvalidLink", func(t *testing.T) {
		db := &mocks.FakeLinkDao{}
		cache := &mocks.FakeLinkDao{}

		r := shortener.NewLinkRepository(db, cache)

		invalid := &shortener.Link{}
		_, err := r.Insert(context.Background(), invalid)

		if !errors.Is(err, shortener.ErrInvalidLink) {
			t.Errorf("Expected err to be %v, but got %v", shortener.ErrInvalidLink, err)
		}
	})

	t.Run("DbFail", func(t *testing.T) {
		db := &mocks.FakeLinkDao{InsertFn: failInsert}
		cache := &mocks.FakeLinkDao{InsertFn: okInsert}

		r := shortener.NewLinkRepository(db, cache)

		link, err := r.Insert(context.Background(), sampleLink)

		if err == nil {
			t.Error("Expected Insert to fail, but got err: nil")
		}

		if cache.InsertCalled {
			t.Error("Expected cache to not have been called")
		}

		if diff := cmp.Diff(sampleLink, link); diff != "" {
			t.Errorf("Inserted link different from expected (-want +got):\n%s", diff)
		}
	})

	t.Run("DbOK", func(t *testing.T) {
		db := &mocks.FakeLinkDao{InsertFn: okInsert}
		cache := &mocks.FakeLinkDao{InsertFn: okInsert}

		r := shortener.NewLinkRepository(db, cache)

		link, err := r.Insert(context.Background(), sampleLink)

		if err != nil {
			t.Errorf("Unexpected error inserting link: %v", err)
		}

		if !cache.InsertCalled {
			t.Error("Expected cache to have been called")
		}

		if diff := cmp.Diff(sampleLink, link); diff != "" {
			t.Errorf("Inserted link different from expected (-want +got):\n%s", diff)
		}
	})
}

func TestDelete(t *testing.T) {
	slug := "b4zoo"
	okDelete := func(ctx context.Context, slug string) error {
		return nil
	}

	failDelete := func(ctx context.Context, slug string) error {
		return shortener.ErrInvalidLink
	}

	t.Run("DbFail", func(t *testing.T) {
		db := &mocks.FakeLinkDao{DeleteFn: failDelete}
		cache := &mocks.FakeLinkDao{DeleteFn: okDelete}

		r := shortener.NewLinkRepository(db, cache)

		err := r.Delete(context.Background(), slug)

		if err == nil {
			t.Error("Expected Insert to fail, but got err: nil")
		}

		if cache.DeleteCalled {
			t.Error("Expected cache to not have been called")
		}
	})

	t.Run("DbOK", func(t *testing.T) {
		db := &mocks.FakeLinkDao{DeleteFn: okDelete}
		cache := &mocks.FakeLinkDao{DeleteFn: okDelete}

		r := shortener.NewLinkRepository(db, cache)

		err := r.Delete(context.Background(), slug)

		if err != nil {
			t.Errorf("Unexpected error inserting link: %v", err)
		}

		if !cache.DeleteCalled {
			t.Error("Expected cache to have been called")
		}
	})
}
