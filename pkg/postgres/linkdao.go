package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

type dao struct {
	conn *pgxpool.Pool
}

// NewLinkDao instantiates a dao for link in postgres db
func NewLinkDao(conn *pgxpool.Pool) shortener.LinkDao {
	return &dao{
		conn: conn,
	}
}

func (d *dao) Find(ctx context.Context, slug string) (*shortener.Link, error) {
	link := shortener.Link{}
	err := d.conn.QueryRow(
		ctx,
		"SELECT url, slug, createdAt FROM links WHERE slug=$1",
		slug,
	).Scan(&link.URL, &link.Slug, &link.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shortener.ErrLinkNotFound
		}
		return nil, err
	}

	return &link, nil
}

func (d *dao) Insert(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
	var createdAt time.Time
	err := d.conn.QueryRow(
		ctx,
		"INSERT INTO links (slug, url) VALUES ($1, $2) RETURNING createdAt",
		l.Slug, l.URL,
	).Scan(&createdAt)

	if err != nil {
		var pgErr *pgconn.PgError
		// if is a unique constraint error code, from postgres
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, shortener.ErrLinkExists
		}
		return nil, err
	}

	l.CreatedAt = createdAt

	return l, nil
}

func (d *dao) Update(ctx context.Context, l *shortener.Link) error {
	panic("not yet implemented")
}

func (d *dao) Delete(ctx context.Context, slug string) error {
	panic("not yet implemented")
}

func (d *dao) List(ctx context.Context, limit int, skip int) ([]shortener.Link, error) {
	rows, err := d.conn.Query(
		ctx,
		`SELECT slug, url, createdAt FROM links LIMIT $1 OFFSET $2`,
		limit,
		skip,
	)

	links := []shortener.Link{}
	if err != nil {
		return links, err
	}

	for rows.Next() {
		l := shortener.Link{}
		err = rows.Scan(&l.Slug, &l.URL, &l.CreatedAt)
		if err != nil {
			break
		}
		links = append(links, l)
	}

	return links, rows.Err()
}
