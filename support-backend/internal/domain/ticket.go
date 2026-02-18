package domain

import "time"

type TicketStatus string //definição do tipo TicketStatus como string

const ( //definição das constantes possíveis para o status do ticket
	TicketOpen    TicketStatus = "OPEN"
	TicketPending TicketStatus = "PENDING"
	TicketClosed  TicketStatus = "CLOSED"
)

type Ticket struct { //definição da struct Ticket: as tags JSON definem como os campos serão serializados/deserializados em JSON
	ID          string       `json:"id"` //`json:"..."` dizem ao encoder JSON qual nome usar ao serializar/deserializar
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TicketStatus `json:"status"`    //json:"status" é do tipo TicketStatus (string com valores pré-definidos)/
	CreatedAt   time.Time    `json:"createdAt"` //time.Time serializa em ISO-8601 por padrão
	UpdatedAt   time.Time    `json:"updatedAt"`
}
