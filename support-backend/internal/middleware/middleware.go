package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler { // Middleware para adicionar um ID único a cada requisição
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Função que será executada para cada requisição
		id := uuid.NewString()             // Gera um novo UUID como ID da requisição
		w.Header().Set("X-Request-Id", id) // Define o cabeçalho de resposta com o ID da requisição
		r.Header.Set("X-Request-Id", id)   // Define o cabeçalho da requisição com o ID da requisição
		next.ServeHTTP(w, r)               // Chama o próximo handler na cadeia
	})
}

func SecurityHeaders(next http.Handler) http.Handler { // Middleware para adicionar cabeçalhos de segurança às respostas HTTP
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Função que será executada para cada requisição
		w.Header().Set("X-Content-Type-Options", "nosniff") // Previne que o navegador faça "sniffing" do tipo de conteúdo
		w.Header().Set("X-Frame-Options", "DENY")           // Impede que a página seja carregada em frames
		w.Header().Set("X-XSS-Protection", "1; mode=block") // Habilita a proteção contra XSS no navegador
		w.Header().Set("Referrer-Policy", "no-referrer")    // Define a política de referência para não enviar o cabeçalho Referer
		next.ServeHTTP(w, r)                                // Chama o próximo handler na cadeia
	})
}

func Logger(next http.Handler) http.Handler { // Middleware para logar informações sobre cada requisição HTTP
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Função que será executada para cada requisição
		start := time.Now()                                                // Registra o tempo de início da requisição
		next.ServeHTTP(w, r)                                               // Chama o próximo handler na cadeia
		log.Printf("%s %s in %s", r.Method, r.URL.Path, time.Since(start)) // Loga o método, caminho e duração da requisição
	})
}

func MaxBodyBytes(limit int64, next http.Handler) http.Handler { // Middleware para limitar o tamanho do corpo da requisição HTTP
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // Função que será executada para cada requisição
		r.Body = http.MaxBytesReader(w, r.Body, limit) // Limita o tamanho do corpo da requisição
		next.ServeHTTP(w, r)                           // Chama o próximo handler na cadeia
	})
}
