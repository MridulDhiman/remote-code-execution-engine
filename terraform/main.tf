provider "aws" {
  region = var.aws_region
}


resource "aws_security_group" "ec2_sg" {
  name        = "ec2_sg"
  description = "Security group for EC2"

  # Inbound Rules
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] ## Allow SSH access from anywhere 
  }

  # Outbound Rules
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "ec2_instance" {
  ami = var.ubuntu_ami_id
  instance_type          = var.instance_type
  key_name               = var.key_name
  vpc_security_group_ids = [aws_security_group.ec2_sg.id]
  user_data = file("hello.sh")

  tags = {
    Name = "ExampleInstance"
  }
}


output "instance_ip" {
  value = aws_instance.ec2_instance.public_ip
}