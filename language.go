package negotiator

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type Lang struct {
	Name    string
	Quality float64
}

// ParseLanguages parses the Accept-Language header and returns a list of languages
// accepted by the client, sorted by priority.
func (n *Negotiator) ParseLanguages(available ...string) ([]string, error) {
	acceptLanguage := n.req.Header.Get("Accept-Language")

	if acceptLanguage == "" || len(available) == 0 || acceptLanguage == "*" {
		return available, nil
	}

	parsedLanguages, err := splitLanguages(acceptLanguage)
	if err != nil {
		return nil, err
	}

	preferredLanguages := findPreferredLanguages(parsedLanguages, available)

	sortLanguagesByPriority(preferredLanguages)

	return getLanguages(preferredLanguages), nil
}

// splitLanguages splits the Accept-Language header into individual languages with quality values.
// It returns an error if a quality value is invalid.
// See https://tools.ietf.org/html/rfc7231#section-5.3.5 for details.
func splitLanguages(acceptLanguage string) ([]Lang, error) {
	languages := strings.Split(acceptLanguage, ",")
	parsedLanguages := make([]Lang, 0, len(languages))

	for _, languageStr := range languages {
		language, err := parseLanguage(languageStr)
		if err != nil {
			return nil, err
		}
		if language.Quality > 0 {
			parsedLanguages = append(parsedLanguages, *language)
		}
	}

	return parsedLanguages, nil
}

// parseLanguage parses a language string into a Lang struct.
// It returns an error if a quality value is invalid.
func parseLanguage(languageStr string) (*Lang, error) {
	language := strings.Split(languageStr, ";q=")
	if len(language) == 1 {
		return &Lang{Name: strings.TrimSpace(language[0]), Quality: 1}, nil
	}

	quality, err := strconv.ParseFloat(strings.TrimSpace(language[1]), 64)
	if err != nil {
		return nil, errors.New("failed to parse quality value")
	}

	return &Lang{Name: strings.TrimSpace(language[0]), Quality: quality}, nil
}

// findPreferredLanguages returns a list of languages that are available.
func findPreferredLanguages(parsedLanguages []Lang, available []string) []Lang {
	availableSet := make(map[string]struct{}, len(available))
	for _, lang := range available {
		availableSet[lang] = struct{}{}
	}

	preferredLanguages := make([]Lang, 0)
	for _, lang := range parsedLanguages {
		if _, ok := availableSet[lang.Name]; ok {
			preferredLanguages = append(preferredLanguages, lang)
		}
	}

	return preferredLanguages
}

// sortLanguagesByPriority sorts a list of languages by priority.
func sortLanguagesByPriority(languages []Lang) {
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Quality > languages[j].Quality
	})
}

func getLanguages(languages []Lang) []string {
	result := make([]string, len(languages))
	for i, language := range languages {
		result[i] = language.Name
	}

	return result
}
