package handler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/joao-fontenele/go-url-shortener/pkg/api/router"
	"github.com/joao-fontenele/go-url-shortener/pkg/mocks"
	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestShortenerRedirect(t *testing.T) {
	linkService := &mocks.FakeLinkService{
		GetURLFn: func(ctx context.Context, slug string) (string, error) {
			if slug == "found" {
				return "https://www.google.com/?search=Google", nil
			}

			if slug == "nFoun" {
				return "", shortener.ErrLinkNotFound
			}

			return "", errors.New("UnexpectedError")
		},
	}
	r := router.New(linkService)

	server := &fasthttp.Server{
		Handler: r.Handler,
	}
	ln := fasthttputil.NewInmemoryListener()

	go server.Serve(ln)
	defer server.Shutdown()

	c := http.Client{
		// don't follow redirects, for the sake of this test case
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		// use custom in memory listener to connect to server
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}
	defer c.CloseIdleConnections()

	tests := []struct {
		Name           string
		Slug           string
		WantBody       []byte
		WantStatusCode int
		WantRedirect   string
	}{
		{
			Name:           "NotFound",
			Slug:           "nFoun",
			WantBody:       []byte(`{"message":"Link with slug 'nFoun' not found","statusCode":404}`),
			WantStatusCode: http.StatusNotFound,
			WantRedirect:   "",
		},
		{
			Name:           "ServerErr",
			Slug:           "error",
			WantBody:       []byte(`{"message":"Error getting slug: UnexpectedError","statusCode":500}`),
			WantStatusCode: http.StatusInternalServerError,
			WantRedirect:   "",
		},
		{
			Name:           "Found",
			Slug:           "found",
			WantBody:       nil,
			WantStatusCode: http.StatusMovedPermanently,
			WantRedirect:   "https://www.google.com/?search=Google",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			endpoint := fmt.Sprintf("http://shortener.com/%s", tc.Slug)
			res, err := c.Get(endpoint)
			if err != nil {
				t.Fatalf("Unexpected error requesting %s: %v", endpoint, err)
			}
			defer res.Body.Close()

			got, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unexpected error parsing response body: %v", err)
			}

			if !bytes.Equal(tc.WantBody, got) {
				t.Errorf("Wrong response (want, got): (%s, %s)", tc.WantBody, got)
			}

			if res.StatusCode != tc.WantStatusCode {
				t.Errorf("Wrong status code (want, got): (%s, %s)", tc.WantBody, got)
			}

			// TODO: find a simpler way to make this assertion
			if tc.WantRedirect != "" {
				if !(len(res.Header["Location"]) == 1 && res.Header["Location"][0] == tc.WantRedirect) {
					t.Errorf("Expected a redirect response, (want, got): (%s, %s)", tc.WantRedirect, res.Header["Location"])
				}
			}
		})
	}
}

func TestNewLink(t *testing.T) {
	linkService := &mocks.FakeLinkService{
		CreateFn: func(ctx context.Context, URL string) (*shortener.Link, error) {
			if URL == "https://ok.com/allOK" {
				return &shortener.Link{
					URL:       "https://ok.com/allOK",
					Slug:      "LolOk",
					CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
				}, nil
			}

			if URL == "" {
				return nil, shortener.ErrInvalidLink
			}

			if URL == "https://link.exists.com" {
				return nil, shortener.ErrLinkExists
			}

			return nil, errors.New("UnexpectedError")
		},
	}
	r := router.New(linkService)

	server := &fasthttp.Server{
		Handler: r.Handler,
	}
	ln := fasthttputil.NewInmemoryListener()

	go server.Serve(ln)
	defer server.Shutdown()

	c := http.Client{
		// use custom in memory listener to connect to server
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}
	defer c.CloseIdleConnections()

	tests := []struct {
		Name           string
		ReqBody        []byte
		WantBody       []byte
		WantStatusCode int
	}{
		{
			Name:           "Ok",
			ReqBody:        []byte(`{"url":"https://ok.com/allOK"}`),
			WantBody:       []byte(`{"slug":"LolOk","url":"https://ok.com/allOK","createdAt":"2020-05-01T00:00:00Z"}`),
			WantStatusCode: http.StatusCreated,
		},
		{
			Name:           "ServerErr",
			ReqBody:        []byte(`{"url":"https://server.error.dev"}`),
			WantBody:       []byte(`{"message":"Error creating link: UnexpectedError","statusCode":500}`),
			WantStatusCode: http.StatusInternalServerError,
		},
		{
			Name:           "InvalidLink",
			ReqBody:        []byte(`{"url":""}`),
			WantBody:       []byte(`{"message":"Link is not valid","statusCode":400}`),
			WantStatusCode: http.StatusBadRequest,
		},
		{
			Name:           "InvalidJSON",
			ReqBody:        []byte(`{"url":"https://link.exists.com",`),
			WantBody:       []byte(`{"message":"Invalid json in request body","statusCode":400}`),
			WantStatusCode: http.StatusBadRequest,
		},
		{
			Name:           "LinkExists",
			ReqBody:        []byte(`{"url":"https://link.exists.com"}`),
			WantBody:       []byte(`{"message":"Error creating link: Link's slug already exists","statusCode":500}`),
			WantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			endpoint := "http://shortener.com/links"
			res, err := c.Post(endpoint, "application/json", bytes.NewReader(tc.ReqBody))
			if err != nil {
				t.Fatalf("Unexpected error requesting %s: %v", endpoint, err)
			}
			defer res.Body.Close()

			got, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unexpected error parsing response body: %v", err)
			}

			if !bytes.Equal(tc.WantBody, got) {
				t.Errorf("Wrong response (want, got): (%s, %s)", tc.WantBody, got)
			}

			if res.StatusCode != tc.WantStatusCode {
				t.Errorf("Wrong status code (want, got): (%d, %d)", tc.WantStatusCode, res.StatusCode)
			}
		})
	}
}

func TestList(t *testing.T) {
	linkService := &mocks.FakeLinkService{
		ListFn: func(ctx context.Context, limit, skip int) ([]shortener.Link, error) {
			links := []shortener.Link{
				{
					URL:       "https://link1.com/?s=google.com",
					Slug:      "link1",
					CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					URL:       "https://link2.com/?s=google.com",
					Slug:      "link2",
					CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					URL:       "https://link3.com/?s=google.com",
					Slug:      "link3",
					CreatedAt: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
				},
			}

			if skip > len(links) {
				return []shortener.Link{}, nil
			}

			if limit+skip > len(links) {
				limit = len(links) - skip
			}

			return links[skip : skip+limit], nil
		},
	}
	r := router.New(linkService)

	server := &fasthttp.Server{
		Handler: r.Handler,
	}
	ln := fasthttputil.NewInmemoryListener()

	go server.Serve(ln)
	defer server.Shutdown()

	c := http.Client{
		// use custom in memory listener to connect to server
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}
	defer c.CloseIdleConnections()

	tests := []struct {
		Name           string
		QueryString    string
		WantBody       []byte
		WantStatusCode int
	}{
		{
			Name:           "LimitNotInt",
			QueryString:    "limit=a&skip=0",
			WantBody:       []byte(`{"message":"Invalid limit argument, must be integer greater than 0","statusCode":400}`),
			WantStatusCode: http.StatusBadRequest,
		},
		{
			Name:           "LimitInvalidInt",
			QueryString:    "limit=0&skip=0",
			WantBody:       []byte(`{"message":"Invalid limit argument, must be integer greater than 0","statusCode":400}`),
			WantStatusCode: http.StatusBadRequest,
		},
		{
			Name:           "SkipNotInt",
			QueryString:    "limit=10&skip=dga",
			WantBody:       []byte(`{"message":"Invalid skip argument, must be integer greater than or equal to 0","statusCode":400}`),
			WantStatusCode: http.StatusBadRequest,
		},
		{
			Name:           "SkpiInvalidInt",
			QueryString:    "limit=10&skip=-1",
			WantBody:       []byte(`{"message":"Invalid skip argument, must be integer greater than or equal to 0","statusCode":400}`),
			WantStatusCode: http.StatusBadRequest,
		},
		{
			Name:           "FirstPageOk",
			QueryString:    "limit=2&skip=0",
			WantBody:       []byte(`[{"slug":"link1","url":"https://link1.com/?s=google.com","createdAt":"2020-05-01T00:00:00Z"},{"slug":"link2","url":"https://link2.com/?s=google.com","createdAt":"2020-05-01T00:00:00Z"}]`),
			WantStatusCode: http.StatusOK,
		},
		{
			Name:           "LastPageOk",
			QueryString:    "limit=2&skip=2",
			WantBody:       []byte(`[{"slug":"link3","url":"https://link3.com/?s=google.com","createdAt":"2020-05-01T00:00:00Z"}]`),
			WantStatusCode: http.StatusOK,
		},
		{
			Name:           "AfterLastPageIsEmpty",
			QueryString:    "limit=2&skip=4",
			WantBody:       []byte(`[]`),
			WantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			endpoint := fmt.Sprintf("http://shortener.com/links?%s", tc.QueryString)
			res, err := c.Get(endpoint)
			if err != nil {
				t.Fatalf("Unexpected error requesting %s: %v", endpoint, err)
			}
			defer res.Body.Close()

			got, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unexpected error parsing response body: %v", err)
			}

			if !bytes.Equal(tc.WantBody, got) {
				t.Errorf("Wrong response (want, got): (%s, %s)", tc.WantBody, got)
			}

			if res.StatusCode != tc.WantStatusCode {
				t.Errorf("Wrong status code (want, got): (%d, %d)", tc.WantStatusCode, res.StatusCode)
			}
		})
	}
}
