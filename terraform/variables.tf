variable "aws_region" {
    type = string
    default = "ap-south-1"
}


variable "instance_type" {
    type = string
    default = "t2.micro"
}

variable "key_name" {
    type = string
    default = "tf-kp"
}

variable "ubuntu_ami_id" {
    type = string
    default = "ami-0f2e255ec956ade7f"
}