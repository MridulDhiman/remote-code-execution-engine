package main

import (
	"gobackend/rabbitmq"
	"gobackend/utils"
	"log"
	"path"

	// "log/slog"
	// "log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	utils.FailOnError(err, "Could not read .env files")
	conn := rabbitmq.NewMQConnection(os.Getenv("AMQP_URL"))

	cwd, err := os.Getwd()
	utils.FailOnError(err, "Could not get current working directory")

	
	// slog.Info("current working directory", "path", cwd + "/");
	arg1 := os.Args[1]
	if arg1 == "worker" {
		if err := conn.Worker(path.Join(cwd,"/docker/javascript/")); err != nil {
			log.Fatal(err)
		}
			
	}
	
	if arg1 == "server" {
		server := NewAPIServer(":8080", conn)
		server.Run()
	}




}
