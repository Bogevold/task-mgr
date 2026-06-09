package task

import (
	"fmt"
)

type InMemoryStore struct {
	tasks  map[uint]Task
	nextID uint
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		tasks: make(map[uint]Task),
	}
}

func (s *InMemoryStore) Save(task Task) (Task, error) {
	task.ID = s.nextID
	s.nextID++
	s.tasks[task.ID] = task
	return task, nil
}

func (s *InMemoryStore) Get(i uint) (Task, error) {
	t, ok := s.tasks[i]
	if ok {
		return t, nil
	}
	return Task{}, fmt.Errorf("task med id %d finnes ikke", i)
}

func (s *InMemoryStore) Delete(i uint) error {
	_, ok := s.tasks[i]
	if ok {
		delete(s.tasks, i)
		return nil
	}
	return fmt.Errorf("kunne ikke slette task med id %d ", i)
}

/*
Sjekk om task.ID finnes i mapen
Hvis ja, erstatt den gamle med den nye
Hvis nei, returner en feil
*/

func (s *InMemoryStore) Update(task Task) (Task, error) {
	_, ok := s.tasks[task.ID]
	if ok {
		s.tasks[task.ID] = task
		return task, nil
	}
	return Task{}, fmt.Errorf("kunne ikke oppdatere task med id %d ", task.ID)
}

func (s *InMemoryStore) GetAll() ([]Task, error) {
	var list []Task
	for _, v := range s.tasks {
		list = append(list, v)
	}
	return list, nil
}

func (s *InMemoryStore) Ping() error {
	return nil
}
