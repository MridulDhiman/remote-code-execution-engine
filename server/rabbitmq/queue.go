package rabbitmq

import (
	"bytes"
	"context"
	"fmt"
	"gobackend/types"
	"gobackend/utils"
	"gobackend/ws"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)


func FailOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }

type MQConnection struct {
	Url string
}

func NewMQConnection(Url string) *MQConnection {
	return &MQConnection{Url: Url}
}

func (mqConn *MQConnection) Init() (*amqp.Connection, *amqp.Channel, amqp.Queue, context.Context, context.CancelFunc) {
	// start connection
	fmt.Println(mqConn.Url)
	conn, err := amqp.Dial(mqConn.Url)
	FailOnError(err, "Could not setup amqp connection")
	// create a channel
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	// create queue
	q, err := ch.QueueDeclare(
		"CodeQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	FailOnError(err, "Could not declare queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return conn, ch, q, ctx, cancel
}

func (mqConn *MQConnection) AddToQueue(input types.CodeExecutionInputBody) error {
    conn, ch, q, ctx, cancel := mqConn.Init()
    defer conn.Close()
	defer ch.Close()
	defer cancel()
	fmt.Println("Initialized connection...")
	inputBuffer, err := utils.EncodeInput(input)
	if err != nil {
		return err
	}
	// adds to queue
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(inputBuffer),
		})
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Successfully added to queue")
	return nil
}



func (mqConn *MQConnection) Worker() error {
   conn, ch, q, _ , cancel := mqConn.Init()
   defer conn.Close()
   defer ch.Close()
   defer cancel()
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	  
	FailOnError(err, "failed to register consumer")

	var forever chan struct{}


	go func() {

		for d := range msgs {
              go ws.Ws.SendMsg("Processing")
			inputBody, err := utils.DecodeInput(d.Body)
			FailOnError(err, "Could not decode the bytes")
			log.Printf("Received a message: %+v", inputBody)
			code := inputBody.Code
			lang:= inputBody.Lang
			input := inputBody.Input

	langConfig := utils.GetLangConfig(lang)
	var inputFileName = "input.txt"
			filePath := path.Join("scripts", langConfig.FileName)
			inputFilePath := path.Join("scripts", inputFileName)
			inputFile, err := os.Create(inputFilePath)
              if err != nil {
				FailOnError(err, "Could not create input file")
			  }
			  defer inputFile.Close()
			  inputFile.WriteString(input)
			file, err := os.Create(filePath)
              if err != nil {
				FailOnError(err, "Could not create file")
			  }
			  defer file.Close()
			  file.WriteString(code)
			  // execute script
			  scriptPath:= "scripts/exec.sh"
			  // make the file executable
			  if err:= os.Chmod(scriptPath, 0700); err!= nil {
				fmt.Println("Error occurred in changing the file mode ", err)
			  }
			  var cmd *exec.Cmd
    // Create the command to execute the shell script based on the OS
    if runtime.GOOS == "windows" {
        // Use Git Bash or WSL on Windows
        bashPath := "C:\\Program Files\\Git\\bin\\bash.exe" // Change this to your Bash path
        cmd = exec.Command(bashPath, scriptPath)
    } else {
        // Use default bash on Unix-like systems
        cmd = exec.Command("/bin/bash", scriptPath)

    }
	_, contName := langConfig.GetImageAndContainerName()

	var fileNameVar = "CODE_FILE=" + langConfig.FileName;
	var inputFileNameVar = "INPUT_FILE=" + inputFileName;
	var contNameVar = "CONTAINER_NAME=" + contName;
	var dockerRunCmdVar = "DOCKER_RUN_CMD=" + langConfig.DockerRunCmd()
	var dockerExecCmdVar = "DOCKER_EXEC_CMD=" + langConfig.DockerExecCmd(len(input) > 0)
	
	fmt.Println(fileNameVar, inputFileNameVar, contNameVar, dockerRunCmdVar, dockerExecCmdVar)
	
	// appends env. variables to command
	
	go ws.Ws.SendMsg("Processed")
	go ws.Ws.SendMsg("Executing")
	cmd.Env = append(os.Environ(), fileNameVar, inputFileNameVar, contNameVar, dockerRunCmdVar, dockerExecCmdVar)
			  // output 
			  output, err := cmd.CombinedOutput()
			  if err != nil {
				fmt.Println("err executing script: ", err)
				return
			  }
			  fmt.Println(string(output))
			go ws.Ws.SendMsg("Finished")
			go ws.Ws.SendMsg(string(output))
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Millisecond)
			log.Printf("Done")
		}
	
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C");
	// it will make the main thread to be forever blocked, unless there is some error
	<-forever

	return nil
}



