#!/bin/bash

CONTAINER_NAME=
IMAGE_NAME=
PEM_FILE=""
REMOTE_USER=""
REMOTE_HOST=""
REMOTE_DIR=""
LOCAL_DIR=""

# Copy the code.js file to the remote server
scp -i "$PEM_FILE" -r $LOCAL_DIR/code.js  ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

# Connect to the remote server and run the necessary Docker commands
ssh -i "$PEM_FILE" ${REMOTE_USER}@${REMOTE_HOST} << EOF
# Run the Docker container with the specified name and image
docker run -dit --name $CONTAINER_NAME  $IMAGE_NAME

# Copy the code.js file from the remote directory into the Docker container
docker cp ${REMOTE_DIR}/code.js $CONTAINER_NAME:/judge/code.js

# Execute the JavaScript file inside the Docker container using Node.js
docker exec $CONTAINER_NAME node /judge/code.js

## Stop and remove the container
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
EOF
