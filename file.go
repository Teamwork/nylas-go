package nylas

// File represents a file in the Nylas system.
type File struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
}
