# Mini Kanban de Tarefas — Desafio Fullstack Veritas

Aplicação de quadro Kanban (React + Go) desenvolvida para a etapa de Desafio
Prático do processo seletivo de estágio Full Stack da **Veritas Consultoria
Empresarial**.

## Stack

- **Frontend:** React 18 + Vite (JavaScript puro, sem TypeScript)
- **Backend:** Go 1.21, apenas biblioteca padrão (`net/http`, `encoding/json`) — sem frameworks externos
- **Persistência:** em memória, com snapshot automático em `backend/tasks.json` a cada alteração (sobrevive a restarts do servidor)

## Estrutura do repositório