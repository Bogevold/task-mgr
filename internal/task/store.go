package task

type TaskStore interface {
	Save(task Task) (Task, error)
	Get(id uint) (Task, error)
	GetAll() ([]Task, error)
	Update(task Task) (Task, error)
	Delete(id uint) error
}
