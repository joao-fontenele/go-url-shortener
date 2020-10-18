package metrics_test

// TODO: refactor these tests to have less code repetition

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/joao-fontenele/go-url-shortener/pkg/metrics"
	"github.com/joao-fontenele/go-url-shortener/pkg/mocks"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
)

func resetMetrics() {
	metrics.HTTPRequestsCounter.Reset()
	metrics.HTTPRequestsDurationHistogram.Reset()
	metrics.DAOFindResultCounter.Reset()
	metrics.DAOOperationsCounter.Reset()
	metrics.DAOOperationsDurationHistogram.Reset()
}

func testMain(m *testing.M) int {
	metrics.Init()
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func TestFind(t *testing.T) {
	link := &shortener.Link{
		URL:       "https://www.google.com",
		Slug:      "aaaaa",
		CreatedAt: time.Now(),
	}
	unexpectedErr := fmt.Errorf("Unexpected")

	successFind := func(ctx context.Context, slug string) (*shortener.Link, error) {
		return link, nil
	}
	errNotFound := func(ctx context.Context, slug string) (*shortener.Link, error) {
		return nil, shortener.ErrLinkNotFound
	}
	UnexpectedErr := func(ctx context.Context, slug string) (*shortener.Link, error) {
		return nil, unexpectedErr
	}

	tests := []struct {
		Name                string
		FindFn              func(ctx context.Context, slug string) (*shortener.Link, error)
		ExpectedLabelResult string
		ExpectedLabelHit    string
		ExpectedErr         error
		ExpectedLink        *shortener.Link
	}{
		{
			Name:                "Success",
			FindFn:              successFind,
			ExpectedLabelResult: "success",
			ExpectedLabelHit:    "hit",
			ExpectedErr:         nil,
			ExpectedLink:        link,
		},
		{
			Name:                "ErrLinkNotFound",
			FindFn:              errNotFound,
			ExpectedLabelResult: "success",
			ExpectedLabelHit:    "miss",
			ExpectedErr:         shortener.ErrLinkNotFound,
			ExpectedLink:        nil,
		},
		{
			Name:                "UnexpectedErr",
			FindFn:              UnexpectedErr,
			ExpectedLabelResult: "error",
			ExpectedLabelHit:    "hit",
			ExpectedErr:         unexpectedErr,
			ExpectedLink:        nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			resetMetrics()

			daoName := "cache"
			baseDao := &mocks.FakeLinkDao{
				FindFn: tc.FindFn,
			}
			dao := metrics.NewLinkDao(baseDao, daoName)

			got, err := dao.Find(context.Background(), "aaaaa")
			if !errors.Is(err, tc.ExpectedErr) {
				t.Errorf("Expected error to be equal %v but got %v", tc.ExpectedErr, err)
			}

			if diff := cmp.Diff(tc.ExpectedLink, got); diff != "" {
				t.Errorf("failed to fetch expected link (-want +got):\n%s", diff)
			}

			metricsText := fmt.Sprintf(
				`
				# HELP dao_find_result_total Total cache hit/miss
				# TYPE dao_find_result_total counter
				dao_find_result_total{name="%s",result="%s"} 1
				`,
				daoName,
				tc.ExpectedLabelHit,
			)
			err = promtest.CollectAndCompare(
				metrics.DAOFindResultCounter,
				strings.NewReader(metricsText),
				"dao_find_result_total",
			)

			if err != nil {
				t.Errorf("Error comparing dao_find_result_total metric: %v", err)
			}

			metricsText = fmt.Sprintf(
				`
				# HELP dao_operations_total Total DAO operations by result (error/success)
				# TYPE dao_operations_total counter
				dao_operations_total{name="%s",result="%s"} 1
				`,
				daoName,
				tc.ExpectedLabelResult,
			)
			err = promtest.CollectAndCompare(
				metrics.DAOOperationsCounter,
				strings.NewReader(metricsText),
				"dao_operations_total",
			)

			if err != nil {
				t.Errorf("Error comparing dao_operations_total metric: %v", err)
			}

			// TODO: test histogram metrics, and find a way to spy on prometheus metrics
			// without having to collect the metrics themselves
		})
	}
}

func TestDelete(t *testing.T) {
	unexpectedErr := fmt.Errorf("Unexpected")

	successDelete := func(ctx context.Context, slug string) error {
		return nil
	}
	errNotFound := func(ctx context.Context, slug string) error {
		return shortener.ErrLinkNotFound
	}
	UnexpectedErr := func(ctx context.Context, slug string) error {
		return unexpectedErr
	}

	tests := []struct {
		Name                string
		DeleteFn            func(ctx context.Context, slug string) error
		ExpectedLabelResult string
		ExpectedErr         error
	}{
		{
			Name:                "Success",
			DeleteFn:            successDelete,
			ExpectedLabelResult: "success",
			ExpectedErr:         nil,
		},
		{
			Name:                "ErrLinkNotFound",
			DeleteFn:            errNotFound,
			ExpectedLabelResult: "success",
			ExpectedErr:         shortener.ErrLinkNotFound,
		},
		{
			Name:                "UnexpectedErr",
			DeleteFn:            UnexpectedErr,
			ExpectedLabelResult: "error",
			ExpectedErr:         unexpectedErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			resetMetrics()

			daoName := "cache"
			baseDao := &mocks.FakeLinkDao{
				DeleteFn: tc.DeleteFn,
			}
			dao := metrics.NewLinkDao(baseDao, daoName)

			err := dao.Delete(context.Background(), "aaaaa")
			if !errors.Is(err, tc.ExpectedErr) {
				t.Errorf("Expected error to be equal %v but got %v", tc.ExpectedErr, err)
			}

			metricsText := fmt.Sprintf(
				`
				# HELP dao_operations_total Total DAO operations by result (error/success)
				# TYPE dao_operations_total counter
				dao_operations_total{name="%s",result="%s"} 1
				`,
				daoName,
				tc.ExpectedLabelResult,
			)
			err = promtest.CollectAndCompare(
				metrics.DAOOperationsCounter,
				strings.NewReader(metricsText),
				"dao_operations_total",
			)

			if err != nil {
				t.Errorf("Error comparing dao_operations_total metric: %v", err)
			}

			// TODO: test histogram metrics, and find a way to spy on prometheus metrics
			// without having to collect the metrics themselves
		})
	}
}

func TestInsert(t *testing.T) {
	link := &shortener.Link{
		URL:       "https://www.google.com",
		Slug:      "aaaaa",
		CreatedAt: time.Now(),
	}
	unexpectedErr := fmt.Errorf("Unexpected")

	successInsert := func(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
		return link, nil
	}
	errNotFound := func(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
		return nil, shortener.ErrLinkNotFound
	}
	UnexpectedErr := func(ctx context.Context, l *shortener.Link) (*shortener.Link, error) {
		return nil, unexpectedErr
	}

	tests := []struct {
		Name                string
		InsertFn            func(ctx context.Context, l *shortener.Link) (*shortener.Link, error)
		ExpectedLabelResult string
		ExpectedErr         error
		ExpectedLink        *shortener.Link
	}{
		{
			Name:                "Success",
			InsertFn:            successInsert,
			ExpectedLabelResult: "success",
			ExpectedErr:         nil,
			ExpectedLink:        link,
		},
		{
			Name:                "ErrLinkNotFound",
			InsertFn:            errNotFound,
			ExpectedLabelResult: "success",
			ExpectedErr:         shortener.ErrLinkNotFound,
			ExpectedLink:        nil,
		},
		{
			Name:                "UnexpectedErr",
			InsertFn:            UnexpectedErr,
			ExpectedLabelResult: "error",
			ExpectedErr:         unexpectedErr,
			ExpectedLink:        nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			resetMetrics()

			daoName := "cache"
			baseDao := &mocks.FakeLinkDao{
				InsertFn: tc.InsertFn,
			}
			dao := metrics.NewLinkDao(baseDao, daoName)

			got, err := dao.Insert(context.Background(), link)
			if !errors.Is(err, tc.ExpectedErr) {
				t.Errorf("Expected error to be equal %v but got %v", tc.ExpectedErr, err)
			}

			if diff := cmp.Diff(tc.ExpectedLink, got); diff != "" {
				t.Errorf("failed to fetch expected link (-want +got):\n%s", diff)
			}

			metricsText := fmt.Sprintf(
				`
				# HELP dao_operations_total Total DAO operations by result (error/success)
				# TYPE dao_operations_total counter
				dao_operations_total{name="%s",result="%s"} 1
				`,
				daoName,
				tc.ExpectedLabelResult,
			)
			err = promtest.CollectAndCompare(
				metrics.DAOOperationsCounter,
				strings.NewReader(metricsText),
				"dao_operations_total",
			)

			if err != nil {
				t.Errorf("Error comparing dao_operations_total metric: %v", err)
			}

			// TODO: test histogram metrics, and find a way to spy on prometheus metrics
			// without having to collect the metrics themselves
		})
	}
}

func TestUpdate(t *testing.T) {
	link := &shortener.Link{
		URL:       "https://www.google.com",
		Slug:      "aaaaa",
		CreatedAt: time.Now(),
	}
	unexpectedErr := fmt.Errorf("Unexpected")

	successUpdate := func(ctx context.Context, l *shortener.Link) error {
		return nil
	}
	errNotFound := func(ctx context.Context, l *shortener.Link) error {
		return shortener.ErrLinkNotFound
	}
	UnexpectedErr := func(ctx context.Context, l *shortener.Link) error {
		return unexpectedErr
	}

	tests := []struct {
		Name                string
		UpdateFn            func(ctx context.Context, l *shortener.Link) error
		ExpectedLabelResult string
		ExpectedErr         error
	}{
		{
			Name:                "Success",
			UpdateFn:            successUpdate,
			ExpectedLabelResult: "success",
			ExpectedErr:         nil,
		},
		{
			Name:                "ErrLinkNotFound",
			UpdateFn:            errNotFound,
			ExpectedLabelResult: "success",
			ExpectedErr:         shortener.ErrLinkNotFound,
		},
		{
			Name:                "UnexpectedErr",
			UpdateFn:            UnexpectedErr,
			ExpectedLabelResult: "error",
			ExpectedErr:         unexpectedErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			resetMetrics()

			daoName := "cache"
			baseDao := &mocks.FakeLinkDao{
				UpdateFn: tc.UpdateFn,
			}
			dao := metrics.NewLinkDao(baseDao, daoName)

			err := dao.Update(context.Background(), link)
			if !errors.Is(err, tc.ExpectedErr) {
				t.Errorf("Expected error to be equal %v but got %v", tc.ExpectedErr, err)
			}

			metricsText := fmt.Sprintf(
				`
				# HELP dao_operations_total Total DAO operations by result (error/success)
				# TYPE dao_operations_total counter
				dao_operations_total{name="%s",result="%s"} 1
				`,
				daoName,
				tc.ExpectedLabelResult,
			)
			err = promtest.CollectAndCompare(
				metrics.DAOOperationsCounter,
				strings.NewReader(metricsText),
				"dao_operations_total",
			)

			if err != nil {
				t.Errorf("Error comparing dao_operations_total metric: %v", err)
			}

			// TODO: test histogram metrics, and find a way to spy on prometheus metrics
			// without having to collect the metrics themselves
		})
	}
}

func TestList(t *testing.T) {
	unexpectedErr := fmt.Errorf("Unexpected")

	successList := func(ctx context.Context, limit, skip int) ([]shortener.Link, error) {
		return []shortener.Link{}, nil
	}
	UnexpectedErr := func(ctx context.Context, limit, skip int) ([]shortener.Link, error) {
		return []shortener.Link{}, unexpectedErr
	}

	tests := []struct {
		Name                string
		ListFn              func(ctx context.Context, limit, skip int) ([]shortener.Link, error)
		ExpectedLabelResult string
		ExpectedErr         error
	}{
		{
			Name:                "Success",
			ListFn:              successList,
			ExpectedLabelResult: "success",
			ExpectedErr:         nil,
		},
		{
			Name:                "UnexpectedErr",
			ListFn:              UnexpectedErr,
			ExpectedLabelResult: "error",
			ExpectedErr:         unexpectedErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			resetMetrics()

			daoName := "cache"
			baseDao := &mocks.FakeLinkDao{
				ListFn: tc.ListFn,
			}
			dao := metrics.NewLinkDao(baseDao, daoName)

			links, err := dao.List(context.Background(), 10, 0)
			if !errors.Is(err, tc.ExpectedErr) {
				t.Errorf("Expected error to be equal %v but got %v", tc.ExpectedErr, err)
			}

			if diff := cmp.Diff([]shortener.Link{}, links); diff != "" {
				t.Errorf("Failed to fetch expected links (-want +got):\n%s", diff)
			}

			metricsText := fmt.Sprintf(
				`
				# HELP dao_operations_total Total DAO operations by result (error/success)
				# TYPE dao_operations_total counter
				dao_operations_total{name="%s",result="%s"} 1
				`,
				daoName,
				tc.ExpectedLabelResult,
			)
			err = promtest.CollectAndCompare(
				metrics.DAOOperationsCounter,
				strings.NewReader(metricsText),
				"dao_operations_total",
			)

			if err != nil {
				t.Errorf("Error comparing dao_operations_total metric: %v", err)
			}

			// TODO: test histogram metrics, and find a way to spy on prometheus metrics
			// without having to collect the metrics themselves
		})
	}
}
