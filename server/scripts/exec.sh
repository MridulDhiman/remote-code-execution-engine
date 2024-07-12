#!/bin/bash

PEM_FILE=""
REMOTE_USER=""
REMOTE_HOST=""
REMOTE_DIR=""
LOCAL_DIR=""


# Copy the code.js file to the remote server
scp -i "$PEM_FILE" -r $LOCAL_DIR/$CODE_FILE  ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}
scp -i "$PEM_FILE" -r $LOCAL_DIR/$INPUT_FILE  ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}

# Connect to the remote server and run the necessary Docker commands
ssh -i "$PEM_FILE" ${REMOTE_USER}@${REMOTE_HOST} << EOF
# Run the Docker container with the specified name and image
$DOCKER_RUN_CMD

# Copy the code.js file from the remote directory into the Docker container
docker cp ${REMOTE_DIR}/$CODE_FILE $CONTAINER_NAME:/judge/$CODE_FILE
docker cp ${REMOTE_DIR}/$INPUT_FILE $CONTAINER_NAME:/judge/$INPUT_FILE

# Execute the code file inside the Docker container using it's corresponding compiler/runtime
$DOCKER_EXEC_CMD

## Stop and remove the container
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
EOF
