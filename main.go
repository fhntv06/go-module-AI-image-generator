package main

import (
	api "AI_image_generator/basic/api"
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
	// important DONT REMOVE / REPLACE
	InitialApp()

	api.ApiKandinsky()
	return

	dbLocal := db.ConstructorDB()
	db.CreateNewTables(dbLocal, "images")

	googlePage := fmt.Sprintf("%s", fetch.Get("https://www.google.com", nil))

	file.CreateDirToPath("google/")
	f := file.OpenFile("google/index.html")
	file.WriteString(googlePage, f)
}
