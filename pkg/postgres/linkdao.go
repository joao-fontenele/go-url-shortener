package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

type dao struct {
	conn *pgxpool.Pool
}

// NewLinkDao instantiates a dao for link in postgres db
func NewLinkDao(conn *pgxpool.Pool) shortener.LinkDao {
	return dao{
		conn: conn,
	}
}

func (d dao) Find(ctx context.Context, slug string) (*shortener.Link, error) {
	link := shortener.Link{}
	err := d.conn.QueryRow(
		ctx,
		"select url, slug, createdAt from links where slug=$1",
		slug,
	).Scan(&link.URL, &link.Slug, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

func (d dao) Insert(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
	if l == nil {
		return nil, fmt.Errorf("invalid link. Link should not be nil")
	}

	var createdAt time.Time
	err := d.conn.QueryRow(
		ctx,
		"INSERT INTO links (slug, url) VALUES ($1, $2) RETURNING createdAt",
		l.Slug, l.URL,
	).Scan(&createdAt)

	if err != nil {
		return nil, err
	}

	l.CreatedAt = createdAt

	return l, nil
}

func (d dao) Update(ctx context.Context, l *shortener.Link) error {
	panic("not yet implemented")
}

func (d dao) Delete(ctx context.Context, slug string) error {
	panic("not yet implemented")
}
