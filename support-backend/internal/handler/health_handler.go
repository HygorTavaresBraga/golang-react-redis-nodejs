package handler

import "net/http" //importa o pacote net/http para manipulação de requisições HTTP

func Health(w http.ResponseWriter, r *http.Request) { //manipulador para checagem de saúde do serviço
	w.WriteHeader(http.StatusOK) //define o status HTTP 200 OK
	w.Write([]byte("OK"))        //escreve "OK" no corpo da resposta
}
