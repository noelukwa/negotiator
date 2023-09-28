package negotiator_test

import (
	"github.com/noelukwa/negotiator"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNegotiator_ParseEncoding(t *testing.T) {

	tests := []struct {
		name           string
		acceptEncoding string
		expected       []string
		available      []string
	}{
		{
			"should return identity",
			"",
			[]string{"identity"},
			nil,
		},
		{
			"should return *",
			"*",
			[]string{"*"},
			nil,
		},
		{
			"should prefer gzip",
			"*, gzip",
			[]string{"gzip", "*"},
			nil,
		},
		{
			"should return *",
			"*, gzip;q=0",
			[]string{"*"},
			nil,
		},
		{
			"should return empty list",
			"gzip;q=0",
			[]string{},
			nil,
		},
		{
			"should return an empty list",
			"*;q=0",
			[]string{},
			nil,
		},
		{
			"should return identity",
			"*;q=0, identity;q=1",
			[]string{"identity"},
			nil,
		},
		{
			"should return identity",
			"identity",
			[]string{"identity"},
			nil,
		},
		{
			"should return an empty list",
			"identity;q=0",
			[]string{},
			nil,
		},
		{
			"should return identity",
			"gzip",
			[]string{"gzip", "identity"},
			nil,
		},
		{
			"should not return compress",
			"gzip, compress;q=0",
			[]string{"gzip", "identity"},
			nil,
		},
		{
			"should return client-preferred encodings",
			"gzip, deflate",
			[]string{"gzip", "deflate", "identity"},
			nil,
		},
		{
			"should return client-preferred encodings",
			"gzip;q=0.8, deflate",
			[]string{"deflate", "gzip", "identity"},
			nil,
		},
		{
			"should return client-preferred encodings",
			"gzip;foo=bar;q=1, deflate;q=1",
			[]string{"gzip", "deflate", "identity"},
			nil,
		},
		{
			"should return client-preferred encodings",
			"gzip;q=0.8, identity;q=0.5, *;q=0.3",
			[]string{"gzip", "identity", "*"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept-Encoding", test.acceptEncoding)

			neg := negotiator.New(req)
			var actual []string

			if test.available == nil {
				actual = neg.ParseEncoding()
			} else {
				actual = neg.ParseEncoding(test.available...)
			}

			if len(test.expected) >= 1 {
				for i, v := range actual {
					if v != test.expected[i] {
						t.Errorf("Expected %s Encoding , got %s", test.expected[i], v)
					}

				}
			} else {
				if len(actual) != len(test.expected) {
					t.Errorf("Expected %s Encodings , got %s", test.expected, actual)
				}
			}

		})
	}

}
