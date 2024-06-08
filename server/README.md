## golang rabbitmq http server

### Local rabbitmq setup

Start Rabbitmq container in detached mode: `docker run -d -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management`

- Port 15672: port to rabbitmq backend dashboard
- Port 5672: port to rabbitmq server


- Create .env file and load AMQP_URL for message queue url.

 Run server: `go run . server`

 Run worker: `go run . worker`

