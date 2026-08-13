# Mini Kanban de Tarefas — Desafio Fullstack Veritas

Aplicação de quadro Kanban (React + Go) desenvolvida para a etapa de Desafio
Prático do processo seletivo de estágio Full Stack da **Veritas Consultoria
Empresarial**.

## Stack

- **Frontend:** React 18 + Vite (JavaScript puro, sem TypeScript)
- **Backend:** Go 1.21, apenas biblioteca padrão (`net/http`, `encoding/json`) — sem frameworks externos
- **Persistência:** em memória, com snapshot automático em `backend/tasks.json` a cada alteração (sobrevive a restarts do servidor)

## Estrutura do repositório

```
desafio-fullstack-veritas/
├── backend/          # API REST em Go
│   ├── main.go
│   ├── handlers.go
│   ├── store.go
│   ├── models.go
│   ├── middleware.go
│   ├── errors.go
│   └── go.mod
├── frontend/          # SPA em React (Vite)
│   └── src/
│       ├── App.jsx
│       ├── api.js
│       ├── constants.js
│       └── components/
│           ├── Column.jsx
│           ├── TaskCard.jsx
│           └── TaskModal.jsx
└── docs/
    └── user-flow.svg  # Diagrama do fluxo de uso do sistema
```

## Como rodar localmente

### 1. Backend (Go)

```bash
cd backend
go run .
```

A API sobe em `http://localhost:8080`. Um arquivo `tasks.json` é criado
automaticamente na primeira escrita, funcionando como persistência simples.

### 2. Frontend (React)

Em outro terminal:

```bash
cd frontend
npm install
cp .env.example .env   # ajuste VITE_API_URL se necessário
npm run dev
```

A aplicação abre em `http://localhost:5173` e já consome a API local.

## Endpoints da API

| Método | Rota          | Descrição                          |
|--------|---------------|-------------------------------------|
| GET    | `/tasks`      | Lista todas as tarefas              |
| GET    | `/tasks/:id`  | Busca uma tarefa específica         |
| POST   | `/tasks`      | Cria uma tarefa                     |
| PUT    | `/tasks/:id`  | Atualiza título, descrição e/ou status |
| DELETE | `/tasks/:id`  | Remove uma tarefa                   |
| GET    | `/health`     | Healthcheck simples                 |

Corpo esperado em `POST`/`PUT`:

```json
{
  "title": "Configurar CORS",
  "description": "Liberar acesso do frontend à API",
  "status": "todo"
}
```

`status` aceita `todo`, `in_progress` ou `done`. No `POST`, se omitido, a
tarefa é criada em `todo`.

## Funcionalidades implementadas

- CRUD completo de tarefas (criar, listar, editar, excluir)
- Três colunas fixas: **A Fazer**, **Em Progresso**, **Concluídas**
- Mover tarefas entre colunas via **drag and drop nativo** (HTML5 Drag and
  Drop API, sem biblioteca externa)
- Feedback visual mínimo: toast de sucesso/erro, estado vazio por coluna,
  loading e tratamento de erro de rede
- Atualização otimista na UI ao mover/excluir tarefas, com rollback em caso
  de falha da API
- CORS liberado no backend para consumo pelo frontend em outra porta

## Decisões técnicas

- **Sem framework/roteador externo no Go:** o escopo é pequeno (uma única
  entidade), então um roteador manual em `handlers.go` mantém o código
  simples e fácil de auditar, sem esconder lógica atrás de "mágica" de
  biblioteca.
- **Persistência híbrida (memória + JSON):** atende ao requisito do edital
  ("armazenamento em memória ou JSON") somando os dois — a fonte da verdade
  em runtime é o mapa em memória (rápido, protegido por mutex), e cada
  escrita é espelhada em disco para não perder dados ao reiniciar o
  processo durante a avaliação.
- **Componentização no frontend:** `Column`, `TaskCard` e `TaskModal` são
  componentes puros e sem estado próprio de dados — todo o estado das
  tarefas vive em `App.jsx`, o que facilita seguir o fluxo de dados durante
  a arguição técnica.
- **Drag and drop nativo:** optou-se pela API nativa do HTML5 em vez de uma
  lib (ex.: `react-beautiful-dnd`) para manter as dependências mínimas e
  deixar explícito o funcionamento do `dataTransfer`.
- **Atualização otimista:** mover ou excluir uma tarefa atualiza a UI
  imediatamente e reverte se a API falhar, priorizando a fluidez pedida no
  critério de avaliação de Funcionalidade.

## Diagrama de User Flow

Ver [`docs/user-flow.svg`](./docs/user-flow.svg) — cobre desde o
carregamento inicial do quadro até os três caminhos possíveis (criar,
mover, editar/excluir) e a convergência de volta ao estado atualizado da
UI.

## Autoria

Projeto desenvolvido individualmente para o desafio prático do processo
seletivo Veritas Consultoria Empresarial.
