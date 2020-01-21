package nylas

// Folder represents a folder in the Nylas system.
type Folder struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
