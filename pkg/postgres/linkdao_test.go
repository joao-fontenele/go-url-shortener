package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joao-fontenele/go-url-shortener/pkg/configger"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

func testMain(m *testing.M) int {
	var err error

	// change dir because default pwd for tests are it's parent dir
	os.Chdir("../../")

	err = configger.Load()
	if err != nil {
		fmt.Printf("failed to load configs: %v", err)
		return 1
	}

	env := configger.Get().Env
	if env != "test" {
		fmt.Println("don't run these tests on non dev environment")
		return 1
	}

	closeDB, err := Connect()
	if err != nil {
		return 1
	}
	defer closeDB()
	return m.Run()
}

func seedDB(conn *pgxpool.Pool) error {
	_, err := conn.Exec(
		context.Background(),
		"INSERT INTO links (slug, url, createdAt) VALUES ('a1CDz', 'https://www.google.com', '2020-05-01T00:00:00.000Z')",
	)

	return err
}

func truncateDB(conn *pgxpool.Pool) error {
	_, err := conn.Exec(context.Background(), "TRUNCATE TABLE links")

	return err
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func TestFind(t *testing.T) {
	conn := GetConnection()
	if err := truncateDB(conn); err != nil {
		t.Fatalf("error truncating test database tables: %v", err)
	}

	err := seedDB(conn)
	if err != nil {
		t.Fatalf("failed to seed db: %v", err)
	}

	dao := NewLinkDao(conn)

	tt := []struct {
		Name string
		Slug string
		Want *shortener.Link
		Err  error
	}{
		{
			Name: "FoundSlug",
			Slug: "a1CDz",
			Want: &shortener.Link{
				URL:       "https://www.google.com",
				Slug:      "a1CDz",
				CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
			},
			Err: nil,
		},
		{
			Name: "NotFoundSlug",
			Slug: "niull",
			Want: nil,
			Err:  shortener.ErrLinkNotFound,
		},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			var got *shortener.Link
			got, err = dao.Find(context.Background(), test.Slug)

			if !errors.Is(err, test.Err) {
				t.Fatalf("failed to find the requested link: %v", err)
			}

			if diff := cmp.Diff(test.Want, got); diff != "" {
				t.Errorf("failed to fetch expected link (-want +got):\n%s", diff)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	conn := GetConnection()
	if err := truncateDB(conn); err != nil {
		t.Fatalf("error truncating test database tables: %v", err)
	}

	err := seedDB(conn)
	if err != nil {
		t.Fatalf("failed to seed db: %v", err)
	}

	dao := NewLinkDao(conn)

	tt := []struct {
		Name  string
		Link  *shortener.Link
		Error error
	}{
		{
			Name: "Success",
			Link: &shortener.Link{
				URL:  "https://www.google.com?s=golang",
				Slug: "spd91",
			},
			Error: nil,
		},
		{
			Name: "ConflictSlug",
			Link: &shortener.Link{
				URL:  "https://www.google.com?s=golang",
				Slug: "a1CDz",
			},
			Error: shortener.ErrLinkExists,
		},
		{
			Name:  "LinkIsNil",
			Link:  nil,
			Error: shortener.ErrInvalidLink,
		},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			_, err := dao.Insert(context.Background(), test.Link)

			if !errors.Is(err, test.Error) {
				t.Fatalf("failed to query new inserted link %v", err)
			}

			if err != nil {
				return
			}

			inserted := shortener.Link{}
			err = conn.QueryRow(
				context.Background(),
				"SELECT slug, url, createdAt FROM links WHERE slug=$1",
				test.Link.Slug,
			).Scan(&inserted.Slug, &inserted.URL, &inserted.CreatedAt)

			if err != nil {
				t.Fatalf("Unexpected error querying inserted link: %v", err)
			}

			if diff := cmp.Diff(test.Link, &inserted); diff != "" {
				t.Errorf("failed to fetch expected link (-want +got):\n%s", diff)
			}
		})
	}
}
