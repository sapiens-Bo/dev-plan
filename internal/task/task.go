package task

type Task struct {
	ID           int64  `json:"id"`
	DeskID       int64  `json:"task_id"`
	Complited    bool   `json:"complited"`
	Description  string `json:"description"`
	ParentTaskID *int64 `json:"parent_task_id,omitempty"`
}

func New(id int64, deskID int64, parentTask *int64, description string) *Task {
	return &Task{
		ID:           id,
		DeskID:       deskID,
		ParentTaskID: parentTask,
		Description:  description,
		Complited:    false,
	}
}

func (t *Task) Done() {
	t.Complited = true
}

func (t *Task) Edit(description string) {
	t.Description = description
}
