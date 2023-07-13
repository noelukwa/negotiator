package negotiator

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type MediaType struct {
	Type       string
	Subtype    string
	Quality    float64
	Parameters map[string]string
}

type Negotiator struct {
	req *http.Request
}

func New(req *http.Request) *Negotiator {
	return &Negotiator{req: req}
}

// ParseMediaTypes parses the Accept header and returns a list of media types
// accepted by the client, sorted by priority.
func (n *Negotiator) ParseMediaTypes(available ...string) []string {
	accept := n.req.Header.Get("Accept")

	if accept == "" {
		accept = "*/*"
	}

	parsedMediaTypes := splitMediaTypes(accept)
	preferredMediaTypes := make([]MediaType, 0)

	for _, mediaType := range parsedMediaTypes {
		if isMediaTypeAccepted(mediaType, available) {
			preferredMediaTypes = append(preferredMediaTypes, mediaType)
		}
	}

	sortMediaTypesByPriority(preferredMediaTypes)

	return getMediaTypes(preferredMediaTypes)
}

// splitMediaTypes splits the Accept header into individual media types with quality values.
func splitMediaTypes(accept string) []MediaType {
	mediaTypes := strings.Split(accept, ",")

	parsedMediaTypes := make([]MediaType, 0)

	for _, mediaTypeStr := range mediaTypes {
		mediaType := parseMediaType(mediaTypeStr)
		if mediaType != nil && mediaType.Quality > 0 {
			parsedMediaTypes = append(parsedMediaTypes, *mediaType)
		}
	}

	return parsedMediaTypes
}

// parseMediaType parses a media type string into a MediaType struct.
func parseMediaType(mediaTypeStr string) *MediaType {
	mediaTypeParts := strings.SplitN(strings.TrimSpace(mediaTypeStr), ";", 2)
	if len(mediaTypeParts) == 0 {
		return nil
	}

	mediaRange := strings.TrimSpace(mediaTypeParts[0])
	mediaTypeParams := make(map[string]string)

	if len(mediaTypeParts) > 1 {
		parameters := splitParameters(mediaTypeParts[1])
		for _, param := range parameters {
			key, val := splitKeyValuePair(param)
			mediaTypeParams[key] = val
		}
	}

	qValue := 1.0

	if qValueStr, ok := mediaTypeParams["q"]; ok {
		q, err := strconv.ParseFloat(qValueStr, 64)
		if err == nil {
			qValue = q
		}
	}

	mediaRangeParts := strings.SplitN(mediaRange, "/", 2)
	if len(mediaRangeParts) != 2 {
		return nil
	}

	return &MediaType{
		Type:       strings.TrimSpace(mediaRangeParts[0]),
		Subtype:    strings.TrimSpace(mediaRangeParts[1]),
		Quality:    qValue,
		Parameters: mediaTypeParams,
	}
}

// isMediaTypeAccepted checks if a media type is accepted by the client.
func isMediaTypeAccepted(mediaType MediaType, available []string) bool {
	if len(available) == 0 {
		return true
	}

	for _, a := range available {
		if matchMediaType(mediaType, a) {
			return true
		}
	}

	return false
}

// matchMediaType checks if a media type matches a specific available media type.
func matchMediaType(mediaType MediaType, available string) bool {
	availableMediaType := parseMediaType(available)
	if availableMediaType == nil {
		return false
	}

	if availableMediaType.Type != "*" && mediaType.Type != availableMediaType.Type {
		return false
	}

	if availableMediaType.Subtype != "*" && mediaType.Subtype != availableMediaType.Subtype {
		return false
	}

	for key, val := range availableMediaType.Parameters {
		if val != "*" && mediaType.Parameters[key] != val {
			return false
		}
	}

	return true
}

// sortMediaTypesByPriority sorts the media types by their priority (q-values).
func sortMediaTypesByPriority(mediaTypes []MediaType) {
	sort.SliceStable(mediaTypes, func(i, j int) bool {
		return mediaTypes[i].Quality > mediaTypes[j].Quality
	})
}

// getMediaTypes returns a list of media types as strings.
func getMediaTypes(mediaTypes []MediaType) []string {
	result := make([]string, len(mediaTypes))

	for i, mediaType := range mediaTypes {
		result[i] = mediaType.Type + "/" + mediaType.Subtype
	}

	return result
}

// splitParameters splits a string of parameters into individual parameter strings.
func splitParameters(paramsStr string) []string {
	parameters := make([]string, 0)

	paramParts := strings.Split(paramsStr, ";")

	for _, param := range paramParts {
		parameters = append(parameters, strings.TrimSpace(param))
	}

	return parameters
}

// splitKeyValuePair splits a key-value pair string into key and value strings.
func splitKeyValuePair(pairStr string) (string, string) {
	pairParts := strings.SplitN(pairStr, "=", 2)

	if len(pairParts) != 2 {
		return "", ""
	}

	key := strings.TrimSpace(pairParts[0])
	value := strings.TrimSpace(pairParts[1])

	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		// Remove quotes from the value
		value = value[1 : len(value)-1]
	}

	return key, value
}
