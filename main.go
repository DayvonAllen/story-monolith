package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"story-app-monolith/database"
	"story-app-monolith/router"
)

func init() {
	database.ConnectToDB()
	database.ConnectToRedis()
}
func main() {
	app := router.Setup()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <- c
		fmt.Println("Shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}
}
