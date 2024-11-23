package desk

type Desk struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func New(id int64, name string) *Desk {
	return &Desk{
		ID:   id,
		Name: name,
	}
}
