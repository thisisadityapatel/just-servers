#!/bin/bash

IMAGE_NAME="thisisadityapatel-just-servers"
TAG="latest"
PORTS=(10000 10001)

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

if ! command_exists docker; then
    echo "Error: Docker not installed."
    exit 1
fi

echo "Building Docker image '$IMAGE_NAME:$TAG'..."
docker build -t "$IMAGE_NAME:$TAG" .
if [ $? -ne 0 ]; then
    echo "Error: Docker build failed."
    exit 1
fi

for port in "${PORTS[@]}"; do
    if lsof -i :"$port" >/dev/null 2>&1; then
        echo "Error: Port $port is already in use."
        exit 1
    fi
done

echo "Running Docker container from '$IMAGE_NAME:$TAG'..."
PORT_ARGS=""
for port in "${PORTS[@]}"; do
    PORT_ARGS="$PORT_ARGS -p $port:$port"
done
docker run -d --name "$IMAGE_NAME" $PORT_ARGS "$IMAGE_NAME:$TAG"
if [ $? -ne 0 ]; then
    echo "Error: Docker run failed."
    exit 1
fi