package negotiator_test

import (
	"github.com/noelukwa/negotiator"
	"net/http/httptest"
	"testing"
)

func TestNegotiator_ParseCharsets(t *testing.T) {

	tests := []struct {
		name          string
		acceptCharset string
		expected      []string
		available     []string
	}{
		{
			"should return *",
			"",
			[]string{"*"},
			nil,
		},
		{
			"should return *",
			"*",
			[]string{"*"},
			nil,
		},
		{
			"should return client-preferred charsets",
			"*, UTF-8",
			[]string{"*", "UTF-8"},
			nil,
		},
		{
			"should exclude UTF-8",
			"*, UTF-8;q=0",
			[]string{"*"},
			nil,
		},
		{
			"should return empty list",
			"UTF-8;q=0",
			[]string{},
			nil,
		},
		{
			"should return client-preferred charsets",
			"ISO-8859-1",
			[]string{"ISO-8859-1"},
			nil,
		},
		{
			"should return client-preferred charsets",
			"UTF-8, ISO-8859-1",
			[]string{"UTF-8", "ISO-8859-1"},
			nil,
		},
		{
			"should return client-preferred charsets",
			"UTF-8;q=0.8, ISO-8859-1",
			[]string{"ISO-8859-1", "UTF-8"},
			nil,
		},
		{
			"should return client-preferred charsets",
			"UTF-8;foo=bar;q=1, ISO-8859-1;q=1",
			[]string{"UTF-8", "ISO-8859-1"},
			nil,
		},
		{
			"should return client-preferred charsets",
			"UTF-8;q=0.9, ISO-8859-1;q=0.8, UTF-8;q=0.7",
			[]string{"UTF-8", "ISO-8859-1"},
			nil,
		},
		{
			"should return empty list for empty list",
			"",
			[]string{},
			[]string{},
		},
		{
			"should return empty list for empty list",
			"*",
			[]string{},
			[]string{},
		},
		{
			"should return matching charsets",
			"*, UTF-8",
			[]string{"UTF-8", "ISO-8859-1"},
			[]string{"UTF-8", "ISO-8859-1"},
		},
		{
			"should exclude UTF-8",
			"*, UTF-8;q=0",
			[]string{"ISO-8859-1"},
			[]string{"UTF-8", "ISO-8859-1"},
		},
		{
			"should always return empty list",
			"UTF-8;q=0",
			[]string{},
			[]string{"UTF-8", "ISO-8859-1"},
		},
		{
			"should return matching charsets",
			"ISO-8859-1",
			[]string{"ISO-8859-1"},
			[]string{"UTF-8", "ISO-8859-1"},
		},
		{
			"should return matching charsets",
			"UTF-8, ISO-8859-1",
			[]string{"UTF-8", "ISO-8859-1"},
			[]string{"UTF-8", "KOI8-R", "ISO-8859-1"},
		},
		{
			"should return matching charsets in client-preferred order",
			"UTF-8;q=0.8, ISO-8859-1",
			[]string{"ISO-8859-1", "UTF-8"},
			[]string{"UTF-8", "KOI8-R", "ISO-8859-1"},
		},
		{
			"should return empty list when no matching charsets",
			"ISO-8859-1",
			[]string{},
			[]string{"utf-8"},
		},
		{
			"should use highest preferred order on duplicate",
			"UTF-8;q=0.9, ISO-8859-1;q=0.8, UTF-8;q=0.7",
			[]string{"UTF-8", "ISO-8859-1"},
			[]string{"ISO-8859-1", "UTF-8"},
		},
		{
			"should return original list",
			"",
			[]string{"UTF-8", "ISO-8859-1"},
			[]string{"UTF-8", "ISO-8859-1"},
		},
		{
			"should be case insensitive, returning provided casing",
			"UTF-8, ISO-8859-1",
			[]string{"UTF-8", "ISO-8859-1"},
			[]string{"utf-8", "iso-8859-1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Accept-Charset", test.acceptCharset)

			neg := negotiator.New(req)
			var actual []string
			if test.available == nil {
				actual = neg.ParseCharsets()
			} else {
				actual = neg.ParseCharsets(test.available...)
			}

			if len(test.expected) >= 1 {
				for i, v := range actual {
					t.Log("actual", v, "expected", test.expected[i], "i", i, "total ex", len(test.expected), "total ac", len(actual), "actual", actual)
					if v != test.expected[i] {
						t.Errorf("Expected %s charset, got %s", test.expected[i], v)
					}
				}
			} else {
				if len(actual) != len(test.expected) {
					t.Errorf("Expected %s charsets, got %s", test.expected, actual)
				}
			}
		})
	}

}
