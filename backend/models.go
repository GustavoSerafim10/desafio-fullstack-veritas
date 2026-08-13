package main

import "time"

// Status representa a coluna em que a tarefa está no quadro Kanban.
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// IsValid confirma se o status recebido é um dos três aceitos pelo sistema.
func (s Status) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	}
	return false
}

// Task representa uma tarefa do quadro Kanban.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TaskInput é o formato aceito no corpo das requisições de criação/edição.
// Campos separados do struct Task evitam que o cliente sobrescreva
// createdAt/updatedAt manualmente.
type TaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

// Validate garante que os dados mínimos exigidos pelo MVP foram enviados.
func (t TaskInput) Validate() error {
	if len(t.Title) == 0 {
		return errEmptyTitle
	}
	if len(t.Title) > 120 {
		return errTitleTooLong
	}
	if t.Status != "" && !t.Status.IsValid() {
		return errInvalidStatus
	}
	return nil
}
