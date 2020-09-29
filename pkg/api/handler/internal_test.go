package handler_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/joao-fontenele/go-url-shortener/pkg/api/router"
	"github.com/joao-fontenele/go-url-shortener/pkg/mocks"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestInternalHandler(t *testing.T) {
	linkService := &mocks.FakeLinkService{}
	r := router.New(linkService)
	server := &fasthttp.Server{
		Handler: r.Handler,
	}

	ln := fasthttputil.NewInmemoryListener()
	go server.Serve(ln)
	defer server.Shutdown()

	c := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}
	defer c.CloseIdleConnections()

	res, err := c.Get("http://shortener.com/internal/status")
	if err != nil {
		t.Fatalf("Unexpected error requesting /internal/status: %v", err)
	}
	defer res.Body.Close()

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Unexpected error parsing response body: %v", err)
	}

	want := []byte(`{"running":true}`)

	if !bytes.Equal(want, got) {
		t.Errorf("Wrong response (want, got): (%s, %s)", want, got)
	}
}
