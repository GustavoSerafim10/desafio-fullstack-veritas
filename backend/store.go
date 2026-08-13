package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// TaskStore mantém as tarefas em memória e persiste o estado em um arquivo
// JSON local a cada alteração, para que os dados sobrevivam a um restart
// do servidor (atende ao requisito "armazenamento em memória ou JSON"
// combinando os dois: fonte da verdade em memória, cópia em disco).
type TaskStore struct {
	mu       sync.RWMutex
	tasks    map[string]Task
	filePath string
	seq      int
}

func NewTaskStore(filePath string) *TaskStore {
	s := &TaskStore{
		tasks:    make(map[string]Task),
		filePath: filePath,
	}
	s.load()
	return s
}

// load lê o arquivo de persistência, se existir. Erros de leitura não são
// fatais: o servidor simplesmente inicia com o quadro vazio.
func (s *TaskStore) load() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return
	}
	for _, t := range tasks {
		s.tasks[t.ID] = t
		s.seq++
	}
}

// persist grava o estado atual no disco. Deve ser chamado com o lock já
// adquirido (write lock) pelo método chamador.
func (s *TaskStore) persist() {
	list := s.allLocked()
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(s.filePath, data, 0644)
}

func (s *TaskStore) nextID() string {
	s.seq++
	return fmt.Sprintf("task-%d", s.seq)
}

func (s *TaskStore) allLocked() []Task {
	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list
}

func (s *TaskStore) All() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.allLocked()
}

func (s *TaskStore) Get(id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, errTaskNotFound
	}
	return t, nil
}

func (s *TaskStore) Create(input TaskInput) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	status := input.Status
	if status == "" {
		status = StatusTodo
	}
	now := time.Now()
	task := Task{
		ID:          s.nextID(),
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[task.ID] = task
	s.persist()
	return task
}

func (s *TaskStore) Update(id string, input TaskInput) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.tasks[id]
	if !ok {
		return Task{}, errTaskNotFound
	}

	existing.Title = input.Title
	existing.Description = input.Description
	if input.Status != "" {
		existing.Status = input.Status
	}
	existing.UpdatedAt = time.Now()

	s.tasks[id] = existing
	s.persist()
	return existing, nil
}

func (s *TaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return errTaskNotFound
	}
	delete(s.tasks, id)
	s.persist()
	return nil
}
