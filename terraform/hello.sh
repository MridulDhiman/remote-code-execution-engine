#!/bin/bash

REMOTE_DIR="/home/ubuntu/code"

sudo apt-get update -y
sudo apt-get install -y docker.io
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ubuntu
sudo apt-get install -y unzip
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

echo "export AWS_DEFAULT_REGION=ap-south-1" >> /etc/environment
mkdir  $REMOTE_DIR
aws s3 sync s3://remotecodeexecutionengine $REMOTE_DIR
docker build -t javascript-image $REMOTE_DIR/javascript
docker build -t gcc-image $REMOTE_DIR/gcc
docker build -t python-image $REMOTE_DIR/python
docker build -t rust-image $REMOTE_DIR/rust
sudo chown ubuntu:ubuntu $REMOTE_DIR