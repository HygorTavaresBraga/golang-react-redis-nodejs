import { apiFetch } from "./http";

export function listTickets() {
  return apiFetch("/api/v1/tickets");
}

export function getTicket(id) {
  return apiFetch(`/api/v1/tickets/${id}`);
}

export function createTicket(payload) {
  return apiFetch("/api/v1/tickets", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function updateTicket(id, payload) {
  return apiFetch(`/api/v1/tickets/${id}`, {
    method: "PATCH",
    body: JSON.stringify(payload),
  });
}

export async function deleteTicket(id) {
  await apiFetch(`/api/v1/tickets/${id}`, { method: "DELETE" });
}
