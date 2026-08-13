export default function Navbar({ apiStatus }) {
  const statusLabel =
    apiStatus === "online"
      ? "API local conectada"
      : apiStatus === "offline"
      ? "API indisponível"
      : "Verificando API...";

  const statusClass =
    apiStatus === "online"
      ? "navbar__status navbar__status--online"
      : apiStatus === "offline"
      ? "navbar__status navbar__status--offline"
      : "navbar__status";

  return (
    <nav className="navbar">
      <div className="navbar__inner">
        <div className="navbar__brand">
          <div className="navbar__logo" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none">
              <rect x="3" y="4" width="6" height="16" rx="1.5" fill="white" fillOpacity="0.95" />
              <rect x="9.5" y="4" width="6" height="10" rx="1.5" fill="white" fillOpacity="0.7" />
              <rect x="16" y="4" width="5" height="6" rx="1.5" fill="white" fillOpacity="0.5" />
            </svg>
          </div>
          <div className="navbar__titles">
            <h1>Mini Kanban de Tarefas</h1>
            <span>Desafio Fullstack · Veritas Consultoria Empresarial</span>
          </div>
        </div>

        <span className={statusClass}>
          <span className="navbar__status-dot" />
          {statusLabel}
        </span>
      </div>
    </nav>
  );
}
