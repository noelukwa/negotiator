package negotiator_test

import (
	"github.com/noelukwa/negotiator"
	"net/http/httptest"
	"testing"
)

func TestNegotiator_ParseCharsets(t *testing.T) {

	casesWithArray := []struct {
		name     string
		header   string
		charsets []string
		expected []string
	}{
		{
			name:     "should return empty list",
			header:   "",
			charsets: []string{},
			expected: []string{},
		},
		{
			name:     "should return original list",
			header:   "",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"UTF-8", "ISO-8859-1"},
		},
		{
			name:     "should return empty list for wildcard",
			header:   "*",
			charsets: []string{},
			expected: []string{},
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
			header:   "*,UTF-8",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"UTF-8", "ISO-8859-1"},
		},
		{
			name:     "should exclude charset with q=0",
			header:   "*, UTF-8;q=0",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"ISO-8859-1"},
		},
		{
			name:     "should return empty list for single q=0 charset",
			header:   "UTF-8;q=0",
			charsets: []string{"UTF-8", "KOI8-R", "ISO-8859-1"},
			expected: []string{},
		},
		{
			name:     "should return matching charsets",
			header:   "ISO-8859-1",
			charsets: []string{"UTF-8", "ISO-8859-1"},
			expected: []string{"ISO-8859-1"},
		},
		{
			name:     "should be case insensitive, returning provided casing in right order",
			header:   "ISO-8859-1",
			charsets: []string{"iso-8859-1", "ISO-8859-1"},
			expected: []string{"iso-8859-1", "ISO-8859-1"},
		},
		{
			name:     "should return empty list when no matching charsets",
			header:   "ISO-8859-1",
			charsets: []string{"utf-8", "KOI8-R"},
			expected: []string{},
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
			charsets: []string{"UTF-8", "ISO-8859-1"},
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
				for i, v := range c.expected {
					if v != actual[i] {
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
