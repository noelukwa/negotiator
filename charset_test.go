package negotiator_test

import (
	"github.com/noelukwa/negotiator"
	"net/http/httptest"
	"testing"
)

func TestNegotiator_ParseCharsets(t *testing.T) {
	cases := []struct {
		name     string
		header   string
		expected []string
	}{
		{
			name:     "should return */*",
			header:   "",
			expected: []string{"*/*"},
		},
		{
			name:   "should return text/*",
			header: "text/*, text/plain;q=0",
			expected: []string{
				"text/*",
			},
		},
		{
			name:   "should return application/json",
			header: "application/json",
			expected: []string{
				"application/json",
			},
		},
		{
			name:   "should return application/json",
			header: "application/json;q=0.5",
			expected: []string{
				"application/json",
			},
		},
		{
			name:     "should return empty slice",
			header:   "application/json;q=0",
			expected: []string{},
		},
		{
			name:   "should return text/html, application/json",
			header: "application/json;q=0.2, text/html",
			expected: []string{
				"text/html",
				"application/json",
			},
		},
		{
			name:   "should return text/*",
			header: "text/*",
			expected: []string{
				"text/*",
			},
		},
		{
			name:   "should return text/plain, text/html, application/json, */*",
			header: "text/plain, application/json;q=0.5, text/html, */*;q=0.1",
			expected: []string{
				"text/plain",
				"text/html",
				"application/json",
				"*/*",
			},
		},
		{
			name:   "should return preferred in order",
			header: "text/plain, application/json;q=0.5, text/html, text/xml, text/yaml, text/javascript, text/csv, text/css, text/rtf, text/markdown, application/octet-stream;q=0.2, */*;q=0.1",
			expected: []string{
				"text/plain",
				"text/html",
				"text/xml",
				"text/yaml",
				"text/javascript",
				"text/csv",
				"text/css",
				"text/rtf",
				"text/markdown",
				"application/json",
				"application/octet-stream",
				"*/*",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Accept", c.header)

			neg := negotiator.New(req)
			actual := neg.ParseMediaTypes()
			if len(c.expected) >= 1 {
				for i, v := range actual {
					if v != c.expected[i] {
						t.Errorf("Expected %s media type, got %s", c.expected[i], v)
					}
				}
			} else {
				if len(actual) != len(c.expected) {
					t.Errorf("Expected %s media types, got %s", c.expected, actual)
				}
			}
		})
	}

	casesWithArray := []struct {
		name     string
		header   string
		charsets []string
		expected []string
	}{
		{
			name:     "should return empty list for empty list",
			header:   "",
			charsets: []string{},
			expected: []string{},
		},
		{
			name:     "should return original list",
			header:   "",
			charsets: []string{"UTF-8"},
			expected: []string{"UTF-8"},
		},
		{
			name:     "should return original list",
			header:   "",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"UTF-8", "ISO-8859-1"},
		},
		{
			name:     "should return empty list for empty list",
			header:   "*",
			charsets: []string{},
			expected: []string{},
		},
		{
			name:     "should return original list",
			header:   "*",
			charsets: []string{"UTF-8"},
			expected: []string{"UTF-8"},
		},
		{
			name:     "should return original list",
			header:   "*",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"UTF-8", "ISO-8859-1"},
		},
		{
			name:     "should return matching charsets",
			header:   "*, UTF-8",
			charsets: []string{"UTF-8"},
			expected: []string{"UTF-8"},
		},
		{
			name:     "should return matching charsets",
			header:   "*, UTF-8",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"UTF-8", "ISO-8859-1"},
		},
		{
			name:     "should exclude UTF-8",
			header:   "*, UTF-8;q=0",
			charsets: []string{"UTF-8"},
			expected: []string{},
		},
		{
			name:     "should exclude UTF-8",
			header:   "*, UTF-8;q=0",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"ISO-8859-1"},
		},
		{
			name:     "should always return empty list",
			header:   "UTF-8;q=0",
			charsets: []string{"ISO-8859-1"},
			expected: []string{},
		},
		{
			name:     "should always return empty list",
			header:   "UTF-8;q=0",
			charsets: []string{"UTF-8", "KOI8-R", "ISO-8859-1"},
			expected: []string{},
		},
		{
			name:     "should always return empty list",
			header:   "UTF-8;q=0",
			charsets: []string{"KOI8-R"},
			expected: []string{},
		},
		{
			name:     "should return matching charsets",
			header:   "ISO-8859-1",
			charsets: []string{"ISO-8859-1"},
			expected: []string{"ISO-8859-1"},
		},
		{
			name:     "should return matching charsets",
			header:   "ISO-8859-1",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"ISO-8859-1"},
		},
		{
			name:     "should be case insensitive, returning provided casing",
			header:   "ISO-8859-1",
			charsets: []string{"iso-8859-1"},
			expected: []string{"iso-8859-1"},
		},
		{
			name:     "should be case insensitive, returning provided casing",
			header:   "ISO-8859-1",
			charsets: []string{"iso-8859-1", "ISO-8859-1"},
			expected: []string{"iso-8859-1", "ISO-8859-1"},
		},
		{
			name:     "should be case insensitive, returning provided casing",
			header:   "ISO-8859-1",
			charsets: []string{"ISO-8859-1", "iso-8859-1"},
			expected: []string{"ISO-8859-1", "iso-8859-1"},
		},
		{
			name:     "should return empty list when no matching charsets",
			header:   "ISO-8859-1",
			charsets: []string{"utf-8"},
			expected: []string{},
		},
		{
			name:     "should return matching charsets",
			header:   "UTF-8, ISO-8859-1",
			charsets: []string{"ISO-8859-1"},
			expected: []string{"ISO-8859-1"},
		},
		{
			name:     "should return matching charsets",
			header:   "UTF-8, ISO-8859-1",
			charsets: []string{"UTF-8", "KOI8-R", "ISO-8859-1"},
			expected: []string{"UTF-8", "ISO-8859-1"},
		},
		{
			name:     "should return empty list when no matching charsets",
			header:   "UTF-8, ISO-8859-1",
			charsets: []string{"KOI8-R"},
			expected: []string{},
		},
		{
			name:     "should return matching charsets in client-preferred order",
			header:   "UTF-8;q=0.8, ISO-8859-1",
			charsets: []string{"ISO-8859-1"},
			expected: []string{"ISO-8859-1"},
		},
		{
			name:     "should return matching charsets in client-preferred order",
			header:   "UTF-8;q=0.8, ISO-8859-1",
			charsets: []string{"UTF-8", "KOI8-R", "ISO-8859-1"},
			expected: []string{"ISO-8859-1", "UTF-8"},
		},
		{
			name:     "should return empty list when no matching charsets",
			header:   "UTF-8;q=0.8, ISO-8859-1",
			charsets: []string{"KOI8-R"},
			expected: []string{},
		},
		{
			name:     "should use highest preferred order on duplicate",
			header:   "UTF-8;q=0.9, ISO-8859-1;q=0.8, UTF-8;q=0.7",
			charsets: []string{"ISO-8859-1"},
			expected: []string{"ISO-8859-1"},
		},
		{
			name:     "should use highest preferred order on duplicate",
			header:   "UTF-8;q=0.9, ISO-8859-1;q=0.8, UTF-8;q=0.7",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"UTF-8", "ISO-8859-1"},
		},
		{
			name:     "should use highest preferred order on duplicate",
			header:   "UTF-8;q=0.9, ISO-8859-1;q=0.8, UTF-8;q=0.7",
			charsets: []string{"ISO-8859-1", "UTF-8"},
			expected: []string{"UTF-8", "ISO-8859-1"},
		},
	}

	for _, c := range casesWithArray {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Accept-Charset", c.header)

			neg := negotiator.New(req)
			actual := neg.ParseCharsets(c.charsets...)
			if len(c.expected) >= 1 {
				for i, v := range actual {
					if v != c.expected[i] {
						t.Errorf("Expected %s charset, got %s", c.expected[i], v)
					}
				}
			} else {
				if len(actual) != len(c.expected) {
					t.Errorf("Expected %s charsets, got %s", c.expected, actual)
				}
			}
		})
	}

}
