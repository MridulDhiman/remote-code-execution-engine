## Code Execution in Remote Containers

### Workflow: 
- [x] You take code, input and language as input, and send HTTP request. 
- [x] It gets added to rabbitmq queue.
- [x] Worker will consume request from from queue and process it inside remote machine with docker pre-installed.
- [x] We will ssh into remote machine and transfer our code and then spin up a new container
- [x] Execute code inside container.
- [x] Return the stderr, stdout to the user.
- [x] Stop and remove the container


### Rabbitmq Cloud  setup: 
1. Login to `https://www.cloudamqp.com/`
2. Create new Rabbitmq instance.
3. Get AMQP instance and add to `.env` file of `/server` directory


### Remote machine Setup:
> Given terraform script configures EC2 Ubuntu instance and docker pre-installed for IAM user.

1. Create IAM user with `EC2FullAccess` and `EC2InstanceConnect` (for SSH) permissions IAM dashboard in AWS console.
2. Install terraform
3. Install AWS CLI
4. Configure aws cli using IAM user access key and secret key: 

```bash
aws configure
```

5. Initialize terraform

```bash
terraform init
```

6. Setup terraform state locally.

```bash
terraform plan -out=FILE
```

7. Apply changes.

```bash
terraform apply
```

8. Destroy the infrastructure

```bash
terraform destroy
```

### Demo:

- Post Request:

![Post Request](public/1.png)

- Output: 

![Output](public/2.png)
