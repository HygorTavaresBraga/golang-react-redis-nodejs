import { useState } from "react";

export default function TicketForm({
  initialTitle = "",
  initialDescription = "",
  initialStatus = "OPEN",
  showStatus = false,
  submitLabel = "Salvar",
  onSubmit,
  disabled = false,
}) {
  const [title, setTitle] = useState(initialTitle);
  const [description, setDescription] = useState(initialDescription);
  const [status, setStatus] = useState(initialStatus);

  function handleSubmit(e) {
    e.preventDefault();
    const payload = { title, description };
    if (showStatus) payload.status = status;
    onSubmit(payload);
  }

  return (
    <form onSubmit={handleSubmit} style={{ display: "grid", gap: 12 }}>
      <label style={{ display: "grid", gap: 6 }}>
        <span>Título</span>
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          minLength={3}
          maxLength={120}
          required
          disabled={disabled}
        />
      </label>

      <label style={{ display: "grid", gap: 6 }}>
        <span>Descrição</span>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          minLength={3}
          maxLength={2000}
          required
          rows={6}
          disabled={disabled}
        />
      </label>

      {showStatus && (
        <label style={{ display: "grid", gap: 6 }}>
          <span>Status</span>
          <select value={status} onChange={(e) => setStatus(e.target.value)} disabled={disabled}>
            <option value="OPEN">OPEN</option>
            <option value="PENDING">PENDING</option>
            <option value="CLOSED">CLOSED</option>
          </select>
        </label>
      )}

      <button className="create-btn" type="submit" disabled={disabled}>
        {submitLabel}
      </button>
    </form>
  );
}
