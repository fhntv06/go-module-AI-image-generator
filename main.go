package main

import (
	api "AI_image_generator/basic/api"
	db "AI_image_generator/basic/database"
	env "AI_image_generator/basic/env_handler"
	"AI_image_generator/basic/fetch"
	"fmt"
	_ "github.com/lib/pq"
	_ "time"
)

func InitialApp() {
	fmt.Println("Initial App!")

	env.InitialEnvParams()
}

func main() {
	InitialApp()

	db.ConstructorDB()
	api.Api()

	fmt.Println(fetch.Get("https://www.google.com"))
}
