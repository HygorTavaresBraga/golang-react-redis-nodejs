package config

import (
	"os"
	"strconv"
)

type Config struct { //define a estrutura Config com as tipagens das variáveis de configuração
	Port      string
	Env       string
	Storage   string
	RedisAddr string
	RedisPass string
	RedisDB   int

	TicketCacheTTLSeconds int
}

func Load() Config { //Load() devolve a struct Config preenchida com as variáveis de ambiente
	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "dev")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPass := getEnv("REDIS_PASSWORD", "")
	ttl := getEnvInt("TICKET_CACHE_TTL_SECONDS", 30)

	return Config{
		Port:      port,
		Env:       env,
		RedisAddr: redisAddr,
		RedisPass: redisPass,
		RedisDB:   0,

		TicketCacheTTLSeconds: ttl,
	}

}

func getEnv(key, fallback string) string { //getEnv retorna o valor da variável de ambiente ou o valor padrão se não estiver definida
	if v := os.Getenv(key); v != "" { //verifica se a variável de ambiente está definida, caso não esteja, retorna String vazia
		return v
	}
	return fallback //retorna o valor padrão se a variável de ambiente não estiver definida
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return fallback
}
