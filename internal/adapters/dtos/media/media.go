package media

type Media struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	MIMEType  string `json:"mimeType"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}
