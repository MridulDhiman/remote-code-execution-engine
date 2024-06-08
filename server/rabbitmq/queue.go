package rabbitmq

import (
	"bytes"
	"context"
	"fmt"
	"gobackend/types"
	"gobackend/utils"
	"runtime"

	// "io"
	"log"
	"os"
	"os/exec"
	"time"

	// docker_types "github.com/docker/docker/api/types"
	// container_types "github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/client"

	// "github.com/docker/docker/pkg/archive"
	// "github.com/docker/go-connections/nat"

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

func (mqConn *MQConnection) Worker(pathDockerImage string) error {
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

			inputBody, err := utils.DecodeInput(d.Body)
			FailOnError(err, "Could not decode the bytes")
			log.Printf("Received a message: %+v", inputBody)
            

			code := inputBody.Code

			file, err := os.Create("scripts/code.js")

              if err != nil {
				FailOnError(err, "Could not create file")
			  }

			  defer file.Close()
			  file.WriteString(code)


			  // execute script
			  scriptPath:= "scripts/exec.sh"
			  // Get and print the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}
	fmt.Println("Current working directory:", wd)

		// Check if the file exists
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			fmt.Println("File does not exist:", scriptPath)
			return
		} else if err != nil {
			fmt.Println("Error checking file:", err)
			return
		}
	

			  // make the file executable
			  if err:= os.Chmod(scriptPath, 0700); err!= nil {
				fmt.Println("Error occurred in changing the file mode ", err)
			  }



			  var cmd *exec.Cmd

    // Create the command to execute the shell script based on the OS
    if runtime.GOOS == "windows" {
        // Use Git Bash or WSL on Windows
        bashPath := "C:\\Program Files\\Git\\bin\\sh.exe" // Change this to your Bash path
        cmd = exec.Command(bashPath, scriptPath)
    } else {
        // Use default bash on Unix-like systems
        cmd = exec.Command("/bin/sh", scriptPath)
    }


			  // output 
			 
			  output, err := cmd.CombinedOutput()
			  if err != nil {
				fmt.Println("err executing script: ", err)
				return
			  }

			  fmt.Println(string(output))

			  
			  

			  //TODO: implement docker build in golang
//            if  err:= os.Setenv("DOCKER_API_VERSION","1.40"); err != nil {
// 			FailOnError(err, "Could not set env. variables")
// 		   }

			
			
// 			apiClient, err := client.NewClientWithOpts(client.FromEnv)
// 			if err != nil {
// 				FailOnError(err, "Docker api client not established")
// 			}


// 		     defer apiClient.Close()


// 			 // archive the folder
// 			// buildCtx, err:= createBuildContext(pathDockerImage)
// 			// if err != nil {
// 			// 	FailOnError(err, "Could not create build context") 
// 			// }

			 

//                ctx:= context.Background()
			   

// 			   imageName:= "javascript-runner"
// 		// resp,  err := apiClient.ImageBuild(ctx, buildCtx, docker_types.ImageBuildOptions{
// 		// 	    Context: buildCtx,
// 		// 		Dockerfile: "Dockerfile",
// 		// 		PullParent: true,
// 		// 		Remove: true,
// 		// 		NoCache: true,
// 		// 		Tags:       []string{imageName},
// 		// 	 })

// 	// 		 cmd := exec.Command("docker", "build", "-t", imageName, "../docker/javascript")
// 	// 		 stdout, err := cmd.Output()

// 	cmd:= exec.Command("docker", "build", "-t", imageName, "../docker/javascript/")
// if err:= cmd.Run(); err!= nil {
// 	FailOnError(err, "Could not build image")
// }


// fmt.Println("Image Built Successfully...")


    

// 			 // create container out of image

// 			 hostBindings:= nat.PortBinding{
// 				HostIP: "0.0.0.0",
// 				HostPort: "8000",
// 			 }

// 			 containerPort, err := nat.NewPort("tcp", "80")
// 	if err != nil {
// 		FailOnError(err, "Unable to get the port")
// 	}



// 	// create file 
	
// 			 cont, err := apiClient.ContainerCreate(
// 				ctx, 
// 				&container_types.Config{
// 				Image: imageName,
// 			 },
// 			 &container_types.HostConfig{
// 				PortBindings: nat.PortMap{containerPort: []nat.PortBinding{hostBindings}},
// 			 },
// 			 nil,
// 			 nil, 
// 			 imageName,
// 			)
			 

// 			if err != nil {
// 				FailOnError(err, "could not create container")
// 			}
// 			apiClient.ContainerStart(ctx, cont.ID, container_types.StartOptions{})

// 			fmt.Printf("Container %s is started\n", cont.ID)
//                // execute code in container 

			    
// 			 res, err :=   apiClient.ContainerLogs(ctx, imageName, container_types.LogsOptions{ 
// 					ShowStdout: true,
// 					 ShowStderr: true,
// 					 Follow: true,
// 					})

                 
// 					if err != nil {
// 						log.Fatal(err)
// 					}


//             //   fmt.Printf("%+v\n", res)
// 			  defer res.Close()
// 			  p := make([]byte, 8)
// 			  res.Read(p)
// 			  content, _ := io.ReadAll(res)
// 			  fmt.Println(string(content))

// 			// apiClient.ContainerStop(ctx, cont.ID, container_types.StopOptions{})
// 			apiClient.ContainerRemove(ctx, cont.ID, container_types.RemoveOptions{
// 				Force: true,
// 			})
           


			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
		}
	
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil

}





// func createBuildContext(path string) (io.Reader, error) {
// 	return archive.TarWithOptions(path, &archive.TarOptions{})
// }