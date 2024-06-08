#!/bin/bash

CONTAINER_NAME=js-cont
IMAGE_NAME=js-image
PEM_FILE=""
REMOTE_USER="ubuntu"
REMOTE_HOST=""
REMOTE_DIR="/home/ubuntu/code"

# Copy the code.js file to the remote server
scp -i "$PEM_FILE" -r code.js ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

# Connect to the remote server and run the necessary Docker commands
ssh -i "$PEM_FILE" ${REMOTE_USER}@${REMOTE_HOST} << EOF
# Run the Docker container with the specified name and image
docker run -d --name $CONTAINER_NAME -v ${REMOTE_DIR}:/jsImage $IMAGE_NAME

# Copy the code.js file from the remote directory into the Docker container
docker cp ${REMOTE_DIR}/code.js $CONTAINER_NAME:/jsImage/code.js

# Execute the JavaScript file inside the Docker container using Node.js
docker exec $CONTAINER_NAME node /jsImage/code.js
EOF
