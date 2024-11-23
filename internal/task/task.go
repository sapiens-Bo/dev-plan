package task

type Task struct {
	ID           int64  `json:"id"`
	DeskID       int64  `json:"task_id"`
	Complited    bool   `json:"complited"`
	Description  string `json:"description"`
	ParentTaskID *int64 `json:"parent_task_id,omitempty"`
}

func New(description string) *Task {
	return &Task{
		Description: description,
	}
}

func (t *Task) Except() {
	t.Complited = true
}

func (t *Task) Edit(description string) {
	t.Description = description
}
