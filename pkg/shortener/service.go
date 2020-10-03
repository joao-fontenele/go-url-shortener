package shortener

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"
)

// slugSize represents the fixed size of a slug
const slugSize = 5

// max attempts for trying to generate a new slug
const maxNewSlugAttempts = 5

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// LinkService will hold the businesses logic to handle link operations
type LinkService interface {
	Create(ctx context.Context, URL string) (*Link, error)
	GetURL(ctx context.Context, slug string) (string, error)
	GetNewSlug(ctx context.Context, size int) (string, error)
	GenerateSlug(size int) string
}

type linkService struct {
	repo LinkRepository
}

// NewLinkService instantiates a LinkService, given a LinkRepository
func NewLinkService(repo LinkRepository) LinkService {
	return &linkService{
		repo: repo,
	}
}

func (ls *linkService) GenerateSlug(size int) string {
	chars := "abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	sb := strings.Builder{}
	sb.Grow(size)
	for i := 0; i < size; i++ {
		sb.WriteByte(chars[seededRand.Intn(len(chars))])
	}
	return sb.String()
}

func (ls *linkService) GetNewSlug(ctx context.Context, size int) (string, error) {
	attempt := 0
	for {
		attempt++
		slug := ls.GenerateSlug(size)

		_, err := ls.repo.Find(ctx, slug)
		if errors.Is(err, ErrLinkNotFound) {
			return slug, nil
		}

		if attempt > maxNewSlugAttempts {
			return "", ErrLinkExists
		}
	}
}

func (ls *linkService) Create(ctx context.Context, URL string) (*Link, error) {
	slug, err := ls.GetNewSlug(ctx, slugSize)

	if err != nil {
		return nil, err
	}

	return ls.repo.Insert(
		ctx,
		&Link{URL: URL, Slug: slug},
	)
}

func (ls *linkService) GetURL(ctx context.Context, slug string) (string, error) {
	l, err := ls.repo.Find(ctx, slug)
	if err != nil {
		return "", err
	}
	return l.URL, err
}
