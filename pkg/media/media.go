package media

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/google/uuid"
)

type Media struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	MIMEType   string `json:"mimeType"`
	Extension  string `json:"extension"`
	Size       int64  `json:"size"`
	TypeBucket Type   `json:"-"`
}

func (m *Media) GetPath() string {
	return filepath.Join(m.TypeBucket.String(), fmt.Sprintf("%s-%s%s", m.ID, m.Title, m.Extension))
}

func (m *Media) GetURLPath() string {
	return filepath.Join(m.TypeBucket.String(), fmt.Sprintf("%s-%s%s", m.ID, url.PathEscape(m.Title), m.Extension))
}

// New - create a new Media object with a unique ID.
func New(title, mimeType, extension string) (*Media, error) {
	typeBucket, err := TypeFromMIME(mimeType)
	if err != nil {
		return nil, err
	}

	return &Media{
		ID:         uuid.New().String(),
		Title:      title,
		TypeBucket: typeBucket,
		MIMEType:   mimeType,
		Extension:  extension,
	}, nil
}
