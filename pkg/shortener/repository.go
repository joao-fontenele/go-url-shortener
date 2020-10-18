package shortener

import (
	"context"
	"fmt"
)

// LinkRepository is a contract between services and underlying datastore
type LinkRepository interface {
	List(ctx context.Context, limit, skip int) ([]Link, error)
	Find(ctx context.Context, slug string) (*Link, error)
	Insert(ctx context.Context, l *Link) (*Link, error)
	Update(ctx context.Context, l *Link) error
	Delete(ctx context.Context, slug string) error
}

type linkRepository struct {
	dbDao    LinkDao
	cacheDao LinkDao
}

var _ LinkRepository = &linkRepository{}

// NewLinkRepository instantiates a LinkRepository, given a Dao
func NewLinkRepository(dbDao LinkDao, cacheDao LinkDao) LinkRepository {
	return &linkRepository{
		dbDao,
		cacheDao,
	}
}

func (lr *linkRepository) Find(ctx context.Context, slug string) (*Link, error) {
	link, _ := lr.cacheDao.Find(ctx, slug)

	if link != nil {
		return link, nil
	}

	link, err := lr.dbDao.Find(ctx, slug)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (lr *linkRepository) Insert(ctx context.Context, l *Link) (*Link, error) {
	err := l.Validate()
	if err != nil {
		return nil, err
	}

	_, err = lr.dbDao.Insert(ctx, l)
	if err != nil {
		return l, err
	}

	// ignore any possible cache errors
	lr.cacheDao.Insert(ctx, l)
	return l, nil
}

func (lr *linkRepository) Update(ctx context.Context, l *Link) error {
	panic("not yet implemented")
}

func (lr *linkRepository) Delete(ctx context.Context, slug string) error {
	err := lr.dbDao.Delete(ctx, slug)
	if err != nil {
		return err
	}

	// ignore any possible cache errors
	lr.cacheDao.Delete(ctx, slug)
	return nil
}

func (lr *linkRepository) List(ctx context.Context, limit int, skip int) ([]Link, error) {
	fmt.Printf("\nrepo: limit=%d, skip=%d\n", limit, skip)
	return lr.dbDao.List(ctx, limit, skip)
}
