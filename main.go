package main

import (
	"flag"
	"log"
	"os"
)

func main() {
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

	startServer()
}
