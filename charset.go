package negotiator

import (
	"sort"
	"strconv"
	"strings"
)

// Charset represents a charset accepted by the client.
type Charset struct {
	Name    string
	Quality float64
}

// ParseCharsets parses the Accept-Charset header and returns a list of charsets
// accepted by the client, sorted by priority.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Charset
func (n *Negotiator) ParseCharsets(available ...string) []string {
	acceptCharset := n.req.Header.Get("Accept-Charset")
	if acceptCharset == "" || acceptCharset == "*" {
		if len(available) == 0 {
			return []string{}
		}
		return available
	}

	parsedCharsets := splitCharsets(acceptCharset)
	uniques := uniqueCharsets(parsedCharsets)

	if len(available) > 0 {
		var filteredCharsets []Charset
		for _, charset := range uniques {
			for _, availCharset := range available {
				if strings.EqualFold(charset.Name, availCharset) {
					filteredCharsets = append(filteredCharsets, charset)
					break
				}
			}
		}
		uniques = filteredCharsets
	}

	sort.SliceStable(uniques, func(i, j int) bool {
		return uniques[i].Quality > uniques[j].Quality
	})

	result := make([]string, 0, len(uniques))
	for _, charset := range uniques {
		if charset.Quality > 0 {
			result = append(result, charset.Name)
		}
	}

	return result
}

// splitCharsets splits the Accept-Charset header into individual charsets with quality values.
func splitCharsets(input string) []Charset {
	rawCharsets := strings.Split(input, ",")
	charsets := make([]Charset, 0, len(rawCharsets))

	for _, rawCharset := range rawCharsets {
		parts := strings.Split(strings.TrimSpace(rawCharset), ";q=")
		charset := Charset{
			Name:    strings.TrimSpace(parts[0]),
			Quality: 1,
		}
		if len(parts) > 1 {
			if quality, err := strconv.ParseFloat(parts[1], 64); err == nil {
				charset.Quality = quality
			}
		}
		charsets = append(charsets, charset)
	}

	return charsets
}

// uniqueCharsets filters the given list of charsets to remove duplicates,
// retaining the highest quality value for each charset name. it then returns
// a new slice of Charset with unique charset names.
func uniqueCharsets(charsets []Charset) []Charset {
	charsetMap := map[string]Charset{}
	for _, charset := range charsets {
		lowerName := strings.ToLower(charset.Name)
		existing, exists := charsetMap[lowerName]
		if !exists || existing.Quality < charset.Quality {
			charsetMap[lowerName] = charset
		}
	}

	unique := make([]Charset, 0, len(charsetMap))
	for _, charset := range charsetMap {
		unique = append(unique, charset)
	}

	return unique
}
