package main

import (
	"log"
	"os"
	"os/signal"

	"saas/env"
	"saas/work"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	workChannel := make(chan string)
	workers, _ := strconv.Atoi(env.Getenv("WORKERS", "3"))
	for i := 0; i < workers; i++ {
		go work.Worker(i, workChannel)
	}
	go work.Dispatcher(workChannel)

	<-stop
	log.Println("\nShutting down the server...")
}
