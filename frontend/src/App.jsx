import { useEffect, useState } from "react";
import Navbar from "./components/Navbar.jsx";
import Column from "./components/Column.jsx";
import TaskModal from "./components/TaskModal.jsx";
import { COLUMNS } from "./constants.js";
import { api } from "./api.js";

function SkeletonBoard() {
  return (
    <div className="skeleton-board" aria-hidden="true">
      {COLUMNS.map((column) => (
        <div className="skeleton-column" key={column.status}>
          <div className="skeleton-line skeleton-title" />
          <div className="skeleton-line skeleton-card" />
          <div className="skeleton-line skeleton-card" style={{ width: "85%" }} />
          <div className="skeleton-line skeleton-card" style={{ width: "70%" }} />
        </div>
      ))}
    </div>
  );
}

export default function App() {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [feedback, setFeedback] = useState(null);
  const [modalState, setModalState] = useState(null); // { mode: 'create'|'edit', status, task }
  const [apiStatus, setApiStatus] = useState("checking"); // 'checking' | 'online' | 'offline'

  useEffect(() => {
    loadTasks();
    pingHealth();
  }, []);

  useEffect(() => {
    if (!feedback) return;
    const timer = setTimeout(() => setFeedback(null), 2500);
    return () => clearTimeout(timer);
  }, [feedback]);

  async function pingHealth() {
    try {
      await api.checkHealth();
      setApiStatus("online");
    } catch {
      setApiStatus("offline");
    }
  }

  async function loadTasks() {
    try {
      setLoading(true);
      setError(null);
      const data = await api.listTasks();
      setTasks(data || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  function openCreateModal(status) {
    setModalState({ mode: "create", status });
  }

  function openEditModal(task) {
    setModalState({ mode: "edit", task });
  }

  async function handleModalSubmit(formValues) {
    try {
      if (modalState.mode === "create") {
        const newTask = await api.createTask({
          ...formValues,
          status: modalState.status,
        });
        setTasks((prev) => [...prev, newTask]);
        setFeedback("Tarefa criada com sucesso.");
      } else {
        const { task } = modalState;
        const updated = await api.updateTask(task.id, {
          ...formValues,
          status: task.status,
        });
        setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
        setFeedback("Tarefa atualizada.");
      }
      setModalState(null);
    } catch (err) {
      setFeedback(`Erro: ${err.message}`);
    }
  }

  async function handleDelete(task) {
    const confirmed = window.confirm(`Excluir a tarefa "${task.title}"?`);
    if (!confirmed) return;

    const previous = tasks;
    setTasks((prev) => prev.filter((t) => t.id !== task.id));
    try {
      await api.deleteTask(task.id);
      setFeedback("Tarefa excluída.");
    } catch (err) {
      setTasks(previous);
      setFeedback(`Erro ao excluir: ${err.message}`);
    }
  }

  async function handleMove(taskId, newStatus) {
    const task = tasks.find((t) => t.id === taskId);
    if (!task || task.status === newStatus) return;

    const previous = tasks;
    setTasks((prev) =>
      prev.map((t) => (t.id === taskId ? { ...t, status: newStatus } : t))
    );
    try {
      await api.updateTask(taskId, {
        title: task.title,
        description: task.description,
        status: newStatus,
      });
    } catch (err) {
      setTasks(previous);
      setFeedback(`Erro ao mover tarefa: ${err.message}`);
    }
  }

  return (
    <>
      <div className="bg-scene" aria-hidden="true">
        <div className="bg-blob bg-blob--1" />
        <div className="bg-blob bg-blob--2" />
        <div className="bg-blob bg-blob--3" />
      </div>

      <Navbar apiStatus={apiStatus} />

      <div className="app">
        {loading && <SkeletonBoard />}

        {!loading && error && (
          <p className="app__status app__status--error">
            Não foi possível carregar as tarefas: {error}.{" "}
            <button onClick={loadTasks}>Tentar novamente</button>
          </p>
        )}

        {feedback && <div className="toast">{feedback}</div>}

        {!loading && !error && (
          <main className="board">
            {COLUMNS.map((column) => (
              <Column
                key={column.status}
                title={column.title}
                status={column.status}
                tasks={tasks.filter((t) => t.status === column.status)}
                onAdd={openCreateModal}
                onEdit={openEditModal}
                onDelete={handleDelete}
                onMove={handleMove}
              />
            ))}
          </main>
        )}

        <footer className="app__footer">
          Mini Kanban · React + Go · desenvolvido para o desafio fullstack Veritas
        </footer>
      </div>

      {modalState && (
        <TaskModal
          initialTask={modalState.task}
          onClose={() => setModalState(null)}
          onSubmit={handleModalSubmit}
        />
      )}
    </>
  );
}
