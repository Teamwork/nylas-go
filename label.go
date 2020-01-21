package nylas

// Label represents a label in the Nylas system.
type Label struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}
