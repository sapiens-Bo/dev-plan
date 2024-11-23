package desk

type Desk struct {
	ID      int
	Name    string
	TasksID []int
}

func (d *Desk) New(name string) *Desk {
	return &Desk{
		Name: name,
	}
}
