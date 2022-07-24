package entities

type Author struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Dob       string `json:"dob,omitempty"`
	PenName   string `json:"pen_name,omitempty"`
}
