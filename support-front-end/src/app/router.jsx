import { createBrowserRouter } from "react-router-dom";
import Layout from "../components/Layout";
import TicketsList from "../pages/TicketsList";
import TicketCreate from "../pages/TicketCreate";
import TicketDetail from "../pages/TicketDetail";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    children: [
      { index: true, element: <TicketsList /> },
      { path: "tickets/new", element: <TicketCreate /> },
      { path: "tickets/:id", element: <TicketDetail /> },
    ],
  },
]);
