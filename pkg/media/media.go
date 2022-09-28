package media

type Media struct {
	ID        string `json:"id" pact:"example=9136ee1d-f0ba-428e-9092-ad64f3ab98863"`
	Title     string `json:"title" pact:"example=Pride and Prejudice"`
	MIMEType  string `json:"mimeType" pact:"example=audio/mp4a-latm"`
	Extension string `json:"extension" pact:"example=.m4b"`
	Size      int64  `json:"size" pact:"example=3869609"`
}

// // New - create a new Media object with a unique ID.
// func New(title, mimeType, extension string) (*Media, error) {
// 	typeBucket, err := TypeFromMIME(mimeType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Media{
// 		ID:         uuid.New().String(),
// 		Title:      title,
// 		MIMEType:   mimeType,
// 		Extension:  extension,
// 	}, nil
// }

// {
// 	"id": "c3bbe250-5532-454f-94cb-47ee0e0588d9",
// 	"title": "Pride and Prejudice",
// 	"mimeType": "audio/mp4a-latm",
// 	"extension": ".m4b",
// 	"size": 3869609
// }
