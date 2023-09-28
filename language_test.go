package negotiator_test

import (
	"github.com/noelukwa/negotiator"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNegotiator_ParseLanguage(t *testing.T) {

	tests := []struct {
		name           string
		acceptLanguage string
		expected       []string
		available      []string
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

			"should return *, en",
			"*, en",
			[]string{"*", "en"},
			nil,
		},
		{

			"should return *",
			"*, en;q=0",
			[]string{"*"},
			nil,
		},
		{

			"should return preferred languages",
			"*;q=0.8, en, es",
			[]string{"en", "es", "*"},
			nil,
		},
		{

			"should return preferred language",
			"en",
			[]string{"en"},
			nil,
		},
		{
			"should return empty list",
			"en;q=0",
			[]string{},
			nil,
		},
		{

			"should return preferred languages",
			"en;q=0.8, es",
			[]string{"es", "en"},
			nil,
		},
		{
			"should use highest preferred order on duplicate",
			"en;q=0.9, es;q=0.8, en;q=0.7",
			[]string{"en", "es"},
			nil,
		},
		{

			"should return en-US, en",
			"en-US, en;q=0.8",
			[]string{"en-US", "en"},
			nil,
		},
		{

			"should return es, en-US",
			"en-US;q=0.8, es",
			[]string{"es", "en-US"},
			nil,
		},
		{

			"should return en-US, en-GB",
			"en-US, en-GB",
			[]string{"en-US", "en-GB"},
			nil,
		},
		{

			"should return es, en-US",
			"en-US;q=0.8, es",
			[]string{"es", "en-US"},
			nil,
		},
		{

			"should return en-US, en-GB",
			"en-US;foo=bar;q=1, en-GB;q=1",
			[]string{"en-US", "en-GB"},
			nil,
		},
		{

			"should use prefer fr over nl",
			"nl;q=0.5, fr, de, en, it, es, pt, no, se, fi, ro",
			[]string{"fr", "de", "en", "it", "es", "pt", "no", "se", "fi", "ro", "nl"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept-Language", test.acceptLanguage)

			neg := negotiator.New(req)
			var actual []string

			if test.available == nil {
				actual, _ = neg.ParseLanguages()
			} else {
				actual, _ = neg.ParseLanguages(test.available...)
			}

			if len(test.expected) >= 1 {
				for i, v := range actual {
					if v != test.expected[i] {
						t.Errorf("Expected %s Language , got %s", test.expected[i], v)
					}
					
				}
			} else {
				if len(actual) != len(test.expected) {
					t.Errorf("Expected %s Languages , got %s", test.expected, actual)
				}
			}

		})
	}
}
