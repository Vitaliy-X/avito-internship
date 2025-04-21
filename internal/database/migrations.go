package database

import (
	"io/ioutil"
	"log"
)

func Migrate() {
	filePath := "scripts/init.sql"

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read %s: %v", filePath, err)
	}

	_, err = DB.Exec(string(content))
	if err != nil {
		log.Fatalf("migration execution error: %v", err)
	}

	log.Println("Migration completed successfully")
}
