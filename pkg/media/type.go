package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Type uint8

const (
	UndefinedType Type = iota
	Audio
	Image
	PDF
)

var (
	errParsingType = errors.New("cannot parse type")
	typeName       = map[uint8]string{
		1: "audio",
		2: "image",
		3: "pdf",
	}
	typeValue = map[string]uint8{
		"audio": 1,
		"mp3":   1,
		"m4a":   1,
		"m4b":   1,
		"image": 2,
		"png":   2,
		"jpg":   2,
		"jpeg":  2,
		"pdf":   3,
	}
)

// String allows Type to implement fmt.Stringer.
func (g Type) String() string {
	return typeName[uint8(g)]
}

// Convert a string to a Type, returns an error if the string is unknown.
func ParseType(s string) (Type, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	value, ok := typeValue[s]
	if !ok {
		return Type(0), fmt.Errorf("%w: %q is not a valid type", errParsingType, s)
	}

	return Type(value), nil
}

// MarshalJSON allows compatibility with marshalling JSON.
func (g Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

// UnmarshalJSON allows compatibility with unmarshalling JSON.
func (g *Type) UnmarshalJSON(data []byte) (err error) {
	var typeName string

	err = json.Unmarshal(data, &typeName)
	if err != nil {
		return err
	}

	if *g, err = ParseType(typeName); err != nil {
		return err
	}

	return nil
}

// TypeFromMIMEType converts a supported MIME type to Type object.
func TypeFromMIME(s string) (Type, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	switch {
	case strings.Contains(s, "audio/"):
		return Audio, nil
	case strings.Contains(s, "image/"):
		return Image, nil
	case strings.Contains(s, "application/pdf"):
		return PDF, nil
	default:
		return Type(0), fmt.Errorf("%w: MIME type %q cannot be converted to Type", errParsingType, s)
	}
}
