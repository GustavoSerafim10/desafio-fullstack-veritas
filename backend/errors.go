package main

import "errors"

var (
	errEmptyTitle    = errors.New("o título da tarefa é obrigatório")
	errTitleTooLong  = errors.New("o título da tarefa deve ter no máximo 120 caracteres")
	errInvalidStatus = errors.New("status inválido: use todo, in_progress ou done")
	errTaskNotFound  = errors.New("tarefa não encontrada")
)
