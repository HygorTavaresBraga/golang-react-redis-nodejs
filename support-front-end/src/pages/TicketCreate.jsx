import { useState } from "react";
import { useNavigate } from "react-router-dom";
import TicketForm from "../components/TicketForm";
import { createTicket } from "../api/tickets";

export default function TicketCreate() {
  const nav = useNavigate();
  const [saving, setSaving] = useState(false);
  const [err, setErr] = useState("");

  async function handleSubmit(payload) {
    setErr("");
    setSaving(true);
    try {
      const t = await createTicket(payload);
      nav(`/tickets/${t.id}`);
    } catch (e) {
      setErr(e.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div style={{ display: "grid", gap: 12 }}>
      <h3 style={{ margin: 0 }}>Novo Ticket</h3>
      {err && <p style={{ color: "crimson" }}>{err}</p>}
      <TicketForm submitLabel={saving ? "Salvando..." : "Criar"} onSubmit={handleSubmit} disabled={saving} />
    </div>
  );
}
