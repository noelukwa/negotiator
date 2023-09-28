package negotiator

// ParseEncoding parses the Accept-Encoding header and returns a list of encodings
// accepted by the client, sorted by priority.
import (
	"sort"
	"strconv"
	"strings"
)

type Encoding struct {
	Name    string
	Quality float64
	Index   int
}

func (n *Negotiator) ParseEncoding(available ...string) []string {
	acceptEncoding := n.req.Header.Get("Accept-Encoding")
	if acceptEncoding == "" {
		return available // If no header is found, return the available encodings as is.
	}

	parsedEncodings := parseAcceptEncoding(acceptEncoding)
	filteredEncodings := filterEncodings(parsedEncodings, available)

	// Sort encodings based on quality and order in the header
	sort.SliceStable(filteredEncodings, func(i, j int) bool {
		if filteredEncodings[i].Quality != filteredEncodings[j].Quality {
			return filteredEncodings[i].Quality > filteredEncodings[j].Quality // Higher quality first
		}
		return filteredEncodings[i].Index < filteredEncodings[j].Index // Original order for same quality
	})

	// Extract the encoding names
	result := make([]string, len(filteredEncodings))
	for i, encoding := range filteredEncodings {
		result[i] = encoding.Name
	}

	return result
}

func parseAcceptEncoding(input string) []Encoding {
	rawEncodings := strings.Split(input, ",")
	encodings := make([]Encoding, 0, len(rawEncodings))

	for i, rawEncoding := range rawEncodings {
		parts := strings.Split(rawEncoding, ";q=")
		name := strings.TrimSpace(parts[0])
		quality := 1.0 // Default quality value
		if len(parts) > 1 {
			q, err := strconv.ParseFloat(parts[1], 64)
			if err == nil {
				quality = q
			}
		}
		encodings = append(encodings, Encoding{Name: name, Quality: quality, Index: i})
	}

	return encodings
}

func filterEncodings(parsedEncodings []Encoding, available []string) []Encoding {
	availableSet := make(map[string]struct{})
	for _, encoding := range available {
		availableSet[strings.ToLower(encoding)] = struct{}{}
	}

	filteredEncodings := make([]Encoding, 0, len(parsedEncodings))
	for _, encoding := range parsedEncodings {
		_, exists := availableSet[strings.ToLower(encoding.Name)]
		if exists {
			filteredEncodings = append(filteredEncodings, encoding)
		}
	}

	return filteredEncodings
}
