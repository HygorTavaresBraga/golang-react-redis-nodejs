import { Link, Outlet } from "react-router-dom";

export default function Layout() {
  return (
    <div style={{ fontFamily: "system-ui", maxWidth: 900, margin: "0 auto", padding: 16 }}>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: 12 }}>
        <Link to="/" style={{ textDecoration: "none", color: "black" }}>
          <h2 style={{ margin: 0 }}>Support - Tickets</h2>
        </Link>
        <nav style={{ display: "flex", gap: 12 }}>
          <Link className="new-btn" to="/tickets/new">Novo</Link>
        </nav>
      </header>

      <hr style={{ margin: "16px 0" }} />

      <Outlet />
    </div>
  );
}
