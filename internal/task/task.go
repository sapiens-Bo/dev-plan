package task

type Task struct {
	ID          int
	Name        string
	IsExcept    bool
	Description string
	UnderTasks  []int
}

func New(name string, description string) *Task {
	return &Task{
		Name:        name,
		Description: description,
	}
}

func (t *Task) Except() {
	t.IsExcept = true
}

func (t *Task) Edit(description string) {
	t.Description = description
}
