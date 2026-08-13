import { useState } from "react";
import TaskCard from "./TaskCard.jsx";

export default function Column({ title, status, tasks, onAdd, onEdit, onDelete, onMove }) {
  const [isDragOver, setIsDragOver] = useState(false);

  function handleDragOver(event) {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
    if (!isDragOver) setIsDragOver(true);
  }

  function handleDragLeave() {
    setIsDragOver(false);
  }

  function handleDrop(event) {
    event.preventDefault();
    setIsDragOver(false);
    const taskId = event.dataTransfer.getData("text/plain");
    if (taskId) onMove(taskId, status);
  }

  return (
    <section
      className={`column column--${status}${isDragOver ? " is-drag-over" : ""}`}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      aria-label={`Coluna ${title}`}
    >
      <header className="column__header">
        <div className="column__title">
          <span className="column__dot" />
          <h2>{title}</h2>
        </div>
        <span className="column__count">{tasks.length}</span>
      </header>

      <div className="column__list">
        {tasks.map((task) => (
          <TaskCard key={task.id} task={task} onEdit={onEdit} onDelete={onDelete} />
        ))}
        {tasks.length === 0 && (
          <p className="column__empty">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" aria-hidden="true">
              <rect x="3" y="7" width="18" height="13" rx="2" />
              <path d="M3 7l3-4h12l3 4" />
              <path d="M9 12h6" />
            </svg>
            <span>Arraste uma tarefa para cá</span>
          </p>
        )}
      </div>

      <button type="button" className="column__add" onClick={() => onAdd(status)}>
        + Nova tarefa
      </button>
    </section>
  );
}
