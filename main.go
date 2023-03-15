package main

import (
	"io/ioutil"
	"log"
	"os"

	"moviesdb.com/config"
	"moviesdb.com/database"
	"moviesdb.com/handler"
)

func main() {
	args := os.Args[1:]
	log.SetOutput(os.Stderr)
	config.Initialize(ioutil.ReadFile)
	conn, err := database.Connect(config.DbConfig())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer conn.Close()

	if len(args) == 1 && args[0] == "migrate" {
		database.RunMigration(conn)
	} else {
		server := handler.NewServer(conn)
		server.Start()
	}
}
