package nylas

// Participant in a Message/Thread.
type Participant struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
