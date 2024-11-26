package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar o arquivo .env: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL não está definido no arquivo .env")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir a conexão com o banco de dados: %v", err)
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			fmt.Println("Conectado ao banco de dados!")
			break
		}
		log.Printf("Tentativa %d: erro ao conectar ao banco: %v. Tentando novamente...\n", i+1, err)
		time.Sleep(time.Second * time.Duration(i*i))
	}

	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco de dados após várias tentativas: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar o driver do banco de dados: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a instância de migração: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("erro ao aplicar migrações: %v", err)
	}

	return db, nil
}
