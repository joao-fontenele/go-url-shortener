package metrics

import (
	"context"
	"errors"
	"time"

	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
	"github.com/prometheus/client_golang/prometheus"
)

type daoWrapper struct {
	dao  shortener.LinkDao
	name string
}

// NewLinkDao returns a dao, that have prometheus metrics calculated automatically
func NewLinkDao(dao shortener.LinkDao, name string) shortener.LinkDao {
	return &daoWrapper{
		name: name,
		dao:  dao,
	}
}

var _ shortener.LinkDao = &daoWrapper{}

func apm(err error, name string, operation string, start time.Time) {
	DAOOperationsDurationHistogram.With(
		prometheus.Labels{"name": name, "operation": operation},
	).Observe(time.Since(start).Seconds())

	result := "error"
	if err == nil || errors.Is(err, shortener.ErrLinkNotFound) {
		result = "success"
	}
	DAOOperationsCounter.With(
		prometheus.Labels{"name": name, "result": result},
	).Inc()
}

func (dw *daoWrapper) Find(ctx context.Context, slug string) (*shortener.Link, error) {
	start := time.Now()
	l, err := dw.dao.Find(ctx, slug)

	findResult := "hit"
	if errors.Is(err, shortener.ErrLinkNotFound) {
		findResult = "miss"
	}
	DAOFindResultCounter.With(
		prometheus.Labels{"name": dw.name, "result": findResult},
	).Inc()

	apm(err, dw.name, "find", start)

	return l, err
}

func (dw *daoWrapper) Insert(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
	l, err := dw.dao.Insert(ctx, l)
	apm(err, dw.name, "insert", time.Now())
	return l, err
}

func (dw *daoWrapper) Update(ctx context.Context, l *shortener.Link) error {
	err := dw.dao.Update(ctx, l)
	apm(err, dw.name, "update", time.Now())
	return err
}

func (dw *daoWrapper) Delete(ctx context.Context, slug string) error {
	err := dw.dao.Delete(ctx, slug)
	apm(err, dw.name, "delete", time.Now())
	return err
}

func (dw *daoWrapper) List(ctx context.Context, limit, skip int) ([]shortener.Link, error) {
	links, err := dw.dao.List(ctx, limit, skip)
	apm(err, dw.name, "list", time.Now())
	return links, err
}
