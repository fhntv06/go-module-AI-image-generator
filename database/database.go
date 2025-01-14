package database

import (
	env "AI_image_generator/basic/env_handler"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DB *sql.DB

var dbLocal *sql.DB

func ConstructorDB() DB {
	connectionParamsDB := GetConnectionParamsDB()

	dbLocal, err := sql.Open("postgres", connectionParamsDB)
	if err != nil {
		log.Fatal("Error in Open:", err)
	}

	fmt.Println("Sql open success!")

	if err := dbLocal.Ping(); err != nil {
		log.Fatal("Error in Ping:", err)
	}

	fmt.Println("Ping success!")

	defer dbLocal.Close()

	return dbLocal
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

func CreateNewTables(name string) {
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

	if _, err := dbLocal.Exec(query); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Таблица", name, "создана.")
}

func CreateUser(username string, password string, createdAt time.Time) int {
	// Выполнение запроса
	var newUserID int
	err := dbLocal.QueryRow(
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

func WorkDB() {
	// test connect PostgreSQL DB
	dbLocal = ConstructorDB()

	fmt.Sprintf("dbLocal %v", dbLocal)
	return

	// sql: database is closed
	CreateNewTables("users")

	newUserID := CreateUser("johndoe23", "secret", time.Now())
	fmt.Printf("Новый пользователь с ID %d успешно добавлен.\n", newUserID)
}
