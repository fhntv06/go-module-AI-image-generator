package database

import (
	env "AI_image_generator/basic/env_handler"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func ConstructorDB() *sql.DB {
	connectionParamsDB := GetConnectionParamsDB()

	instanceDB, err := sql.Open("postgres", connectionParamsDB)
	if err != nil {
		log.Fatal("Error in Open:", err)
	}

	fmt.Println("Sql open success!")

	if err := instanceDB.Ping(); err != nil {
		log.Fatal("Error in Ping:", err)
	}

	fmt.Println("Ping success!")

	return instanceDB
}
func CloseDB(instanceDB *sql.DB) {
	defer instanceDB.Close()
}

func GetConnectionParamsDB() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.GetEnvParam("APP_DB_HOST"),
		env.GetEnvParam("APP_DB_PORT"),
		env.GetEnvParam("APP_DB_USER"),
		env.GetEnvParam("APP_DB_PASSWORD"),
		env.GetEnvParam("APP_DB_DATABASE"),
	)
}

func CreateNewTables(instanceDB *sql.DB, name string) {
	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		name,
	)

	if _, err := instanceDB.Exec(query); err != nil {
		fmt.Println("Таблица ", name, "не создана.")
		log.Fatal(err)
	}

	fmt.Println("Таблица", name, "создана.")

	CloseDB(instanceDB)
}

func CreateUser(instanceDB *sql.DB, username string, password string, createdAt time.Time) int {
	// Выполнение запроса
	var newUserID int
	err := instanceDB.QueryRow(
		`INSERT INTO $1 (username, email, password_hash) VALUES ($2, $3, $4) RETURNING id`,
		env.GetEnvParam("DATABASE"),
		username,
		password,
		createdAt,
	).Scan(&newUserID)

	if err != nil {
		log.Fatal(err)
	}

	return newUserID
}
