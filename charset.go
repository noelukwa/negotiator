package negotiator

import (
	"fmt"
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
func (n *Negotiator) ParseCharsets(available ...string) []string {
	acceptCharset := n.req.Header.Get("Accept-Charset")

	fmt.Println("acceptCharset:", acceptCharset)
	if acceptCharset == "" || len(available) == 0 || acceptCharset == "*" {
		return available
	}

	parsedCharsets := splitCharsets(acceptCharset)
	preferredCharsets := make([]Charset, 0)
	fmt.Println("parsedCharsets:", parsedCharsets)
	for _, charset := range parsedCharsets {

		if idx, ok := isCharsetAccepted(charset, available); ok {
			charset.Name = available[idx]
			preferredCharsets = append(preferredCharsets, charset)
		}
	}

	fmt.Println("preferredCharsets:", preferredCharsets)
	sortCharsetsByPriority(preferredCharsets)

	return getCharsets(preferredCharsets)
}

// splitCharsets splits the Accept-Charset header into individual charsets with quality values.
func splitCharsets(acceptCharset string) []Charset {
	charsets := strings.Split(acceptCharset, ",")

	parsedCharsets := make([]Charset, 0)

	for _, charsetStr := range charsets {
		charset := parseCharset(charsetStr)
		if charset != nil && charset.Quality > 0 {
			parsedCharsets = append(parsedCharsets, *charset)
		}
	}

	return parsedCharsets
}

// parseCharset parses a single charset string with its quality value from the Accept-Charset header.
func parseCharset(charsetStr string) *Charset {
	charsetParts := strings.SplitN(strings.TrimSpace(charsetStr), ";", 2)
	if len(charsetParts) == 0 {
		return nil
	}

	charsetName := strings.TrimSpace(charsetParts[0])

	qValue := 1.0

	if len(charsetParts) > 1 {
		parameters := splitParameters(charsetParts[1])
		for _, param := range parameters {
			key, val := splitKeyValuePair(param)
			if key == "q" {
				q, err := strconv.ParseFloat(val, 64)
				if err == nil {
					qValue = q
				}
			}
		}
	}

	return &Charset{
		Name:    charsetName,
		Quality: qValue,
	}
}

// isCharsetAccepted checks if a charset is accepted by the client.
func isCharsetAccepted(charset Charset, available []string) (int, bool) {
	if len(available) == 0 {
		return -1, true
	}

	for i, a := range available {
		fmt.Printf("a: %s, charset.Name: %s\n", a, charset.Name)
		if strings.EqualFold(a, charset.Name) {
			return i, true
		}
	}

	return -1, false
}

// sortCharsetsByPriority sorts the charsets by their priority (quality value).
func sortCharsetsByPriority(charsets []Charset) {
	sort.SliceStable(charsets, func(i, j int) bool {
		return charsets[i].Quality > charsets[j].Quality

	})
}

// getCharsets returns a list of charsets as strings.
func getCharsets(charsets []Charset) []string {
	result := make([]string, len(charsets))

	for i, charset := range charsets {
		result[i] = charset.Name
	}

	return result
}
