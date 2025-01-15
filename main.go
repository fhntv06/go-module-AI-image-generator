package main

import (
	db "AI_image_generator/basic/database"
	env "AI_image_generator/basic/env_handler"
	fetch "AI_image_generator/basic/fetch"
	file "AI_image_generator/basic/file"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "time"
)

var dbLocal *sql.DB

func InitialApp() {
	fmt.Println("Initial App!")

	env.InitialEnvParams()
}

func main() {
	InitialApp()

	dbLocal := db.ConstructorDB()
	db.CreateNewTables(dbLocal, "images")

	googlePage := fmt.Sprintf("%s", fetch.Get("https://www.google.com"))

	file.CreateDirToPath("google/")
	f := file.OpenFile("google/index.html")
	file.WriteString(googlePage, f)
}
