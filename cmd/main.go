package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/victorlui/sma-api/internal/db"
	"github.com/victorlui/sma-api/internal/delivery/http"
)

func main() {
	// inicialização do banco local
	// conn := "postgres://sma-admin:sma-admin@localhost:5432/sma-db?sslmode=disable"
	// postgresDB, err := db.NewConnection(conn)

	// Inicializa o banco de dados DOCKER
	dsn := os.Getenv("DATABASE_URL")
	postgresDB, err := db.NewConnection(dsn)

	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer postgresDB.Close()

	r := http.SetupRouter(postgresDB)

	//Finaliza o servidor e a conexão em caso de sinal de encerramento
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Encerrando servidor...")
		postgresDB.Close()
		os.Exit(0)
	}()

	// Inicia o servidor
	log.Println("Servidor rodando na porta 9090")
	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
