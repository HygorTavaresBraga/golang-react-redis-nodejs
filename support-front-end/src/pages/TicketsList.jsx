import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listTickets, deleteTicket } from "../api/tickets";
import { trash, pencilSharp } from "ionicons/icons";

export default function TicketsList() {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState("");
  const statusColors = {
    OPEN: "green",
    PENDING: "orange",
    CLOSED: "red",
  };

  async function refresh() {
    setErr("");
    setLoading(true);
    try {
      const data = await listTickets();
      setItems(data);
    } catch (e) {
      setErr(e.message);
    } finally {
      setLoading(false);
    }
  }

 useEffect(() => {
  refresh();

  console.log("[SSE] connecting...");
  const es = new EventSource("/api/v1/tickets/stream");

  es.onopen = () => console.log("[SSE] open");

  // Se você manteve "event: ticket" no backend:
  es.addEventListener("ticket", () => {
    console.log("[SSE] ticket event -> refresh");
    refresh();
  });

  // Se você mudar o backend para enviar só "data:" (message padrão), isso funciona:
  es.onmessage = () => {
    console.log("[SSE] message -> refresh");
    refresh();
  };

  es.onerror = (e) => {
    console.log("[SSE] error", e);
  };

  return () => {
    console.log("[SSE] close");
    es.close();
  };
  // eslint-disable-next-line react-hooks/exhaustive-deps
}, []);


  async function handleDelete(id) {
    if (!window.confirm("Deletar este ticket?")) return;
    try {
      await deleteTicket(id);
      await refresh();
    } catch (e) {
      alert(e.message);
    }
  }

  if (loading) return <p>Carregando...</p>;

  return (
    <div style={{ display: "grid", gap: 12 }}>

      {err && <p style={{ color: "crimson" }}>{err}</p>}

      {items.length === 0 ? (
        <p>Nenhum ticket ainda.</p>
      ) : (
        <table border="1" cellPadding="8" style={{ borderCollapse: "collapse", width: "100%" }}>
          <thead className="thead-list">
            <tr>
              <th>ID</th>
              <th>Título</th>
              <th>Status</th>
              <th>Criado</th>
              <th>Ações</th>
            </tr>
          </thead>
          <tbody className="tbody-list">
            {items.map((t) => (
              <tr key={t.id}>
                <td style={{ fontFamily: "monospace" }}>{t.id.slice(0, 8)}...</td>
                <td>
                  <Link to={`/tickets/${t.id}`}>{t.title}</Link>
                </td>

                <td style={{ color: statusColors[t.status] || "inherit" }}>
                  {t.status}
                </td>

                <td>{new Date(t.createdAt).toLocaleString()}</td>
                
                <td style={{ display: "flex", gap: 8, justifyContent: "center" }}>
                  <Link className="edit-btn" to={`/tickets/${t.id}`}>
                    <img src={pencilSharp} alt="Editar" width="18" height="18" />
                  </Link>
                  <button className="delete-btn" onClick={() => handleDelete(t.id)} aria-label="Deletar">
                    <img src={trash} alt="Deletar" width="18" height="18" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
