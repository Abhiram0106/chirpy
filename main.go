package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if secretsErr := godotenv.Load(); secretsErr != nil {
		log.Fatal(secretsErr)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		log.Println("USING DEBUG MODE")
		log.Println("Deleting database.json...")
		err := os.Remove(databasePath)
		if err != nil {
			log.Println(err)
		}
	}

	cfg := &apiConfig{
		fileserverHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	startServer(cfg)
}
