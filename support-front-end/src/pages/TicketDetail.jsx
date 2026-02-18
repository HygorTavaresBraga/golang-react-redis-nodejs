import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { deleteTicket, getTicket, updateTicket } from "../api/tickets";
import TicketForm from "../components/TicketForm";

export default function TicketDetail() {
  const { id } = useParams();
  const nav = useNavigate();
  const [ticket, setTicket] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [err, setErr] = useState("");

  async function load() {
    setErr("");
    setLoading(true);
    try {
      const t = await getTicket(id);
      setTicket(t);
    } catch (e) {
      setErr(e.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, [id]);

  async function handleSubmit(payload) {
    setErr("");
    setSaving(true);
    try {
      const updated = await updateTicket(id, payload);
      setTicket(updated);
      alert("Atualizado!");
    } catch (e) {
      setErr(e.message);
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete() {
    if (!window.confirm("Deletar este ticket?")) return;
    try {
      await deleteTicket(id);
      nav("/");
    } catch (e) {
      alert(e.message);
    }
  }

  if (loading) return <p>Carregando...</p>;
  if (err) return <p style={{ color: "crimson" }}>{err}</p>;
  if (!ticket) return <p>Ticket não encontrado.</p>;

  return (
    <div style={{ display: "grid", gap: 12 }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: 12 }}>
        <div>
          <h3 style={{ margin: 0 }}>Ticket</h3>
          <div style={{ fontFamily: "monospace", opacity: 0.8 }}>{ticket.id}</div>
        </div>
        <button onClick={handleDelete}>Deletar</button>
      </div>

      <TicketForm
        initialTitle={ticket.title}
        initialDescription={ticket.description}
        initialStatus={ticket.status}
        showStatus={true}
        submitLabel={saving ? "Salvando..." : "Salvar alterações"}
        onSubmit={handleSubmit}
        disabled={saving}
      />

      <div style={{ display: "grid", gap: 6, fontSize: 14, opacity: 0.8 }}>
        <div>
          Status atual: <b>{ticket.status}</b>
        </div>
        <div>Criado: {new Date(ticket.createdAt).toLocaleString()}</div>
        <div>Atualizado: {new Date(ticket.updatedAt).toLocaleString()}</div>
      </div>
    </div>
  );
}
