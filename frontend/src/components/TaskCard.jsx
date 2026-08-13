import { useState } from "react";

export default function TaskCard({ task, onEdit, onDelete }) {
  const [isDragging, setIsDragging] = useState(false);

  function handleDragStart(event) {
    event.dataTransfer.setData("text/plain", task.id);
    event.dataTransfer.effectAllowed = "move";
    setIsDragging(true);
  }

  function handleDragEnd() {
    setIsDragging(false);
  }

  return (
    <article
      className={`task-card${isDragging ? " is-dragging" : ""}`}
      draggable
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <header className="task-card__header">
        <h3>{task.title}</h3>
        <div className="task-card__actions">
          <button
            type="button"
            title="Editar tarefa"
            onClick={() => onEdit(task)}
          >
            ✏️
          </button>
          <button
            type="button"
            title="Excluir tarefa"
            onClick={() => onDelete(task)}
          >
            🗑️
          </button>
        </div>
      </header>
      {task.description && <p className="task-card__desc">{task.description}</p>}
    </article>
  );
}
