package shortener_test

import (
	"errors"
	"testing"
	"time"

	"github.com/joao-fontenele/go-url-shortener/pkg/shortener"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		Name    string
		Input   *shortener.Link
		WantErr error
	}{
		{
			Name: "ValidLink",
			Input: &shortener.Link{
				Slug:      "aaaaa",
				CreatedAt: time.Now(),
				URL:       "https://www.google.com?search=Google#frag",
			},
			WantErr: nil,
		},
		{
			Name:    "InvalidNilLink",
			Input:   nil,
			WantErr: shortener.ErrInvalidLink,
		},
		{
			Name: "InvalidEmptyURL",
			Input: &shortener.Link{
				Slug:      "aaaaa",
				CreatedAt: time.Now(),
				URL:       "",
			},
			WantErr: shortener.ErrInvalidLink,
		},
		{
			Name: "InvalidNoSchemeURL",
			Input: &shortener.Link{
				Slug:      "aaaaa",
				CreatedAt: time.Now(),
				URL:       "google.com",
			},
			WantErr: shortener.ErrInvalidLink,
		},
		{
			Name: "InvalidNoHostRL",
			Input: &shortener.Link{
				Slug:      "aaaaa",
				CreatedAt: time.Now(),
				URL:       "http://",
			},
			WantErr: shortener.ErrInvalidLink,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			got := tc.Input.Validate()
			if tc.WantErr == nil && got != nil {
				t.Errorf("Expected error to be nil, but got %#v", got)
			} else if !errors.Is(got, tc.WantErr) {
				t.Errorf("Expected error to be %#v but got %#v", tc.WantErr, got)
			}
		})
	}
}
