package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type TaskHandler struct {
	store *TaskStore
}

func NewTaskHandler(store *TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

// ServeHTTP funciona como um roteador simples para o prefixo /tasks,
// evitando a dependência de um router externo para um escopo tão pequeno.
func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/tasks")
	path = strings.Trim(path, "/")

	switch {
	case path == "" && r.Method == http.MethodGet:
		h.list(w, r)
	case path == "" && r.Method == http.MethodPost:
		h.create(w, r)
	case path != "" && r.Method == http.MethodPut:
		h.update(w, r, path)
	case path != "" && r.Method == http.MethodDelete:
		h.delete(w, r, path)
	case path != "" && r.Method == http.MethodGet:
		h.getOne(w, r, path)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método não suportado para esta rota")
	}
}

func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.All())
}

func (h *TaskHandler) getOne(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	var input TaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	if err := input.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	task := h.store.Create(input)
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var input TaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}
	if err := input.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	task, err := h.store.Update(id, input)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
