package main

import (
	"avito-internship/internal/app"
	"avito-internship/internal/database"
)

func main() {
	database.Connect()
	database.Migrate()
	app.Run()
}
